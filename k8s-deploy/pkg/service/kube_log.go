/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service/kube"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

// GetLogs is Get container Logs
func GetLogs(args *LogChanArgs) (*bytes.Buffer, error) {

	log.Debug().Interface("args", args).Msg("GetLogs")

	var logReadList = make(map[string]io.ReadCloser)

	list, err := getLogRead(args)
	if err != nil {
		return nil, err
	}

	logReadList = list

	log.Debug().Int("len", len(logReadList)).Msg("buf")
	buf := new(bytes.Buffer)

	var prefix func(string) string
	if len(logReadList) == 1 {
		prefix = func(s string) string {
			return ""
		}
	} else {
		prefix = func(s string) string {
			return s
		}
	}

	w := bufio.NewWriter(buf)

	defer w.Flush()
	for k, v := range logReadList {
		log.Debug().Str("k", k).Msg("for")
		msg, err := readLogToString(v, prefix("["+k+"] "))
		if err != nil {
			return nil, err
		}
		log.Debug().Str("k", k).Msg("for1")
		w.WriteString(msg)
	}

	return buf, nil

}

func readLogToString(logRead io.ReadCloser, prefix string) (string, error) {
	defer logRead.Close()
	defer log.Debug().Str("prefix", prefix).Msg("readLogToQueue close")
	buf := new(bytes.Buffer)
	r := bufio.NewReader(logRead)
	log.Debug().Msg("for")
	for {
		msgstr, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Warn().Err(err).Msg("read log form io read error")
				return "", err
			}
			log.Debug().Err(err).Msg("read log form io read io.EOF")
			break
		}
		buf.WriteString(prefix + msgstr)
	}
	return buf.String(), nil
}

// getLogFollow is Get container Logs
func getLogRead(args *LogChanArgs) (map[string]io.ReadCloser, error) {

	podContainerList, err := getPodContainerList(args.Name, getDefaultNamespace(args.Namespace), args.Container)
	if err != nil {
		return nil, err
	}

	readCloserList := make(map[string]io.ReadCloser)

	for containerName, podName := range podContainerList {
		readCloser, err := KubeClient.GetPodLogs(getDefaultNamespace(args.Namespace), podName, &kube.PodLogArgs{
			Container:                    containerName,
			Follow:                       args.Follow,
			Previous:                     args.Previous,
			SinceSeconds:                 args.SinceSeconds,
			SinceTime:                    args.SinceTime,
			Timestamps:                   args.Timestamps,
			TailLines:                    args.TailLines,
			LimitBytes:                   args.LimitBytes,
			InsecureSkipTLSVerifyBackend: args.InsecureSkipTLSVerifyBackend,
		})
		if err != nil {
			return nil, err
		}
		key := getLogPrefix(containerName, podName)
		readCloserList[key] = readCloser
		log.Debug().Str("key", key).Interface("v", readCloser).Msg("got io podContainerList ")
	}

	return readCloserList, nil
}

func getLogPrefix(containerName, podName string) string {
	return fmt.Sprintf("%s %s", podName, containerName)
}

type LogChanArgs struct {
	Name                         string
	Namespace                    string
	Container                    string
	Follow                       bool
	Previous                     bool
	SinceSeconds                 *int64
	SinceTime                    time.Time
	Timestamps                   bool
	TailLines                    *int64
	LimitBytes                   *int64
	InsecureSkipTLSVerifyBackend bool
}

// WriteLog WriteLog
func WriteLog(w *websocket.Conn, args *LogChanArgs) (err error) {
	defer w.Close()
	log.Debug().Interface("args", args).Msg("WriteLog")
	queue := utils.NewQueue(128)

	var logReadList = make(map[string]io.ReadCloser)

	list, err := getLogRead(args)
	if err != nil {
		return err
	}

	logReadList = list

	defer func() {
		for k, v := range logReadList {
			v.Close()
			log.Debug().Str("key", k).Msg("io readCloser close")
		}
	}()

	var prefix func(string) string
	if len(logReadList) == 1 {
		prefix = func(s string) string {
			return ""
		}
	} else {
		prefix = func(s string) string {
			return s
		}
	}

	for k, v := range logReadList {
		go readLogToQueue(v, prefix("["+k+"] "), queue)
	}

	for {
		v, ok, _ := queue.Get()
		if !ok {
			if err := websocket.Message.Send(w, ""); err != nil {
				return err
			}
		} else {
			if err := websocket.Message.Send(w, fmt.Sprint(v)); err != nil {
				return err
			}
		}
	}
}

func readLogToQueue(logRead io.ReadCloser, prefix string, queue utils.Queue) error {
	defer logRead.Close()
	defer log.Debug().Str("prefix", prefix).Msg("readLogToQueue close")
	r := bufio.NewReader(logRead)
	for {
		msgstr, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Warn().Err(err).Msg("read log form io read error")
				return err
			}
			log.Debug().Err(err).Msg("read log form io read io.EOF")
			return nil
		}

		for queue.Quantity() > queue.Capacity() {
			time.Sleep(time.Millisecond)
		}
		queue.Put(fmt.Sprintf("%s%s", prefix, msgstr))
	}
}

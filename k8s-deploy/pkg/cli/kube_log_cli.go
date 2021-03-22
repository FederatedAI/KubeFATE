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

package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/websocket"
)

func LogCommand() *cli.Command {
	return &cli.Command{
		Name: "logs",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "follow",
				Aliases: []string{"f"},
				Value:   false,
				Usage:   "Specify if the logs should be streamed.",
			},
			&cli.BoolFlag{
				Name:  "previous",
				Value: false,
				Usage: "If true, print the logs for the previous instance of the container in a pod if it exists.",
			},
			&cli.DurationFlag{
				Name:  "since",
				Usage: "Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time since may be used.",
			},
			&cli.TimestampFlag{
				Name:   "since-time",
				Layout: "2006-01-02T15:04:05",
				Usage:  "Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time since may be used.",
			},
			&cli.BoolFlag{
				Name:  "timestamps",
				Usage: "Include timestamps on each line in the log output.",
			},
			&cli.Int64Flag{
				Name:  "tail",
				Value: -1,
				Usage: "Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided.",
			},
			&cli.Int64Flag{
				Name:  "limit-bytes",
				Value: 0,
				Usage: "Maximum bytes of logs to return. Defaults to no limit.",
			},
		},
		Usage: "Get this Cluster module log",
		Action: func(c *cli.Context) error {

			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}

			args := url.Values{}
			args.Set("container", fmt.Sprint(c.Args().Get(1)))
			if c.Bool("previous") {
				args.Set("previous", "true")
			}
			if c.Duration("since") != 0 {
				args.Set("since", fmt.Sprint(c.Duration("since").Seconds()))
			}

			if c.Timestamp("since-time") != nil {
				args.Set("since-time", fmt.Sprint(c.Timestamp("since-time").Format(time.RFC3339)))
			}

			if c.Bool("timestamps") {
				args.Set("timestamps", "true")
			}
			if c.Int64("tail") != -1 {
				args.Set("tail", fmt.Sprint(c.Int64("tail")))
			}
			if c.Int64("limit-bytes") != 0 {
				args.Set("limit-bytes", fmt.Sprint(c.Int64("limit-bytes")))
			}

			log.Debug().Str("args", args.Encode()).Msg("args.Encode")

			follow := c.Bool("follow")

			if follow {
				return GetModuleLogFollow(uuid, args.Encode())
			}

			kubeLog, err := GetModuleLog(uuid, args.Encode())
			if err != nil {
				return err
			}
			fmt.Println(kubeLog)
			return nil
		},
	}
}

func GetModuleLog(uuid, args string) (string, error) {
	r := &Request{
		Type: "GET",
		Path: "log",
		Body: nil,
	}

	serviceUrl := viper.GetString("serviceurl")
	apiVersion := api.APIVersion + "/"
	if serviceUrl == "" {
		serviceUrl = "localhost:8080/"
	}
	Url := "http://" + serviceUrl + "/" + apiVersion + r.Path + fmt.Sprintf("/%s?%s", uuid, args)

	body := bytes.NewReader(r.Body)
	log.Debug().Str("Type", r.Type).Str("url", Url).Msg("Request")
	request, err := http.NewRequest(r.Type, Url, body)
	if err != nil {
		return "", err
	}

	token, err := getToken()
	if err != nil {
		return "", err
	}
	Authorization := fmt.Sprintf("Bearer %s", token)

	request.Header.Add("Authorization", Authorization)
	request.Header.Add("user-agent", "kubefate")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Debug().Int("StatusCode", resp.StatusCode).Interface("body", resp.Body).Msg("resp Status Code")

	if resp.StatusCode != 200 {
		type LogResultErr struct {
			Error string
		}

		logResultErr := new(LogResultErr)

		err = json.Unmarshal(respBody, &logResultErr)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("resp.StatusCode=%d, error: %s", resp.StatusCode, logResultErr.Error)
	}

	type LogResultMsg struct {
		Msg  string
		Data string
	}

	LogResult := new(LogResultMsg)

	err = json.Unmarshal(respBody, &LogResult)
	if err != nil {
		return "", err
	}

	log.Debug().Int("Code", resp.StatusCode).Msg("ok")
	return LogResult.Data, err
}

func GetModuleLogFollow(uuid, args string) error {

	r := &Request{
		Type: "GET",
		Path: "log",
		Body: nil,
	}

	serviceUrl := viper.GetString("serviceurl")
	apiVersion := api.APIVersion + "/"
	if serviceUrl == "" {
		serviceUrl = "localhost:8080/"
	}
	Url := "ws://" + serviceUrl + "/" + apiVersion + r.Path + fmt.Sprintf("/%s/ws?%s", uuid, args)
	log.Debug().Str("Url", Url).Msg("ok")

	config, err := websocket.NewConfig(Url, "http://"+serviceUrl+"/")
	config.Header.Add("user-agent", "kubefate")
	token, err := getToken()
	if err != nil {
		return err
	}
	Authorization := fmt.Sprintf("Bearer %s", token)
	config.Header.Add("Authorization", Authorization)
	ws, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		var msg string
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			if err != io.EOF {
				log.Err(err).Msg("Receive form websocket error")
				return err
			}
			log.Debug().Err(err).Msg("Receive io.EOF")
			return nil
		}
		fmt.Printf(msg)
	}
}

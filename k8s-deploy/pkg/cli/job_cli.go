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
	"io/ioutil"
	"net/http"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func JobCommand() *cli.Command {
	return &cli.Command{
		Name:  "job",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			JobListCommand(),
			JobInfoCommand(),
			JobStopCommand(),
			JobDeleteCommand(),
		},
		Usage: "List jobs, describe and delete a job",
	}
}

func JobListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags:   []cli.Flag{},
		Usage:   "Show job list",
		Action: func(c *cli.Context) error {
			job := new(Job)
			return GetItemList(job)
		},
	}
}

func JobDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Flags:   []cli.Flag{},
		Usage:   "Delete a job",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			job := new(Job)
			return DeleteItem(job, uuid)
		},
	}
}

func JobInfoCommand() *cli.Command {
	return &cli.Command{
		Name: "describe",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "uuid",
				Value: "",
				Usage: "Describe a job with given UUID",
			},
		},
		Usage: "Show job's details info",
		Action: func(c *cli.Context) error {

			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			job := new(Job)
			return GetItem(job, uuid)
		},
	}
}

func JobStopCommand() *cli.Command {
	return &cli.Command{
		Name: "stop",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "uuid",
				Value: "",
				Usage: "Describe a job with given UUID",
			},
		},
		Usage: "Stop job",
		Action: func(c *cli.Context) error {

			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}

			r := &Request{
				Type: "PUT",
				Path: "job",
				Body: nil,
			}

			serviceURL := viper.GetString("serviceurl")
			apiVersion := api.APIVersion + "/"
			if serviceURL == "" {
				serviceURL = "localhost:8080/"
			}
			URL := "http://" + serviceURL + "/" + apiVersion + r.Path + "/" + uuid + "?jobStatus=stop"
			body := bytes.NewReader(r.Body)
			log.Debug().Str("Type", r.Type).Str("url", URL).Msg("Request")
			request, err := http.NewRequest(r.Type, URL, body)
			if err != nil {
				return err
			}

			token, err := getToken()
			if err != nil {
				return err
			}
			Authorization := fmt.Sprintf("Bearer %s", token)

			request.Header.Add("Authorization", Authorization)
			request.Header.Add("user-agent", "kubefate")
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if resp.StatusCode != 200 {
				type JobErrMsg struct {
					Error string
				}
				jobErrMsg := new(JobErrMsg)
				err = json.Unmarshal(respBody, &jobErrMsg)
				if err != nil {
					return err
				}
				return fmt.Errorf("resp.StatusCode=%d, error: %s", resp.StatusCode, jobErrMsg.Error)
			}

			type JobResultMsg struct {
				Msg  string
				Data string
			}

			JobResult := new(JobResultMsg)

			err = json.Unmarshal(respBody, &JobResult)
			if err != nil {
				return err
			}

			log.Debug().Int("Code", resp.StatusCode).Bytes("Body", respBody).Msg("ok")

			fmt.Println(JobResult.Data)
			return nil

		},
	}
}

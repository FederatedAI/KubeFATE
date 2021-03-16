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
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"io/ioutil"
)

func ChartCommand() *cli.Command {
	return &cli.Command{
		Name:  "chart",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			ChartListCommand(),
			ChartInfoCommand(),
			ChartDeleteCommand(),
			ChartCreateCommand(),
		},
		Usage: "List Charts, create, delete and describe a Chart",
	}
}

func ChartListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags:   []cli.Flag{},
		Usage:   "List Charts list",
		Action: func(c *cli.Context) error {
			chart := new(Chart)
			return GetItemList(chart)
		},
	}
}

func ChartInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "describe",
		Flags: []cli.Flag{},
		Usage: "Describe a Chart's detail info",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			chart := new(Chart)
			return GetItem(chart, uuid)
		},
	}
}

func ChartDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Flags:   []cli.Flag{},
		Usage:   "Delete a Chart",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}

			chart := new(Chart)
			log.Debug().Str("uuid", uuid).Msg("Chart delete uuid")
			return DeleteItem(chart, uuid)

		},
	}
}

func ChartCreateCommand() *cli.Command {

	return &cli.Command{
		Name: "upload",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "Upload a Chart with local given file",
			},
		},
		Usage: "Upload a Chart from local",
		Action: func(c *cli.Context) error {

			file := c.String("file")

			log.Debug().Str("file", file).Msg("file")

			filename := filepath.Base(file)
			log.Debug().Str("filename", filename).Msg("filename")

			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)

			//
			fileWriter, err := bodyWriter.CreateFormFile("file", filename)
			if err != nil {
				fmt.Println("error writing to buffer")
				return err
			}

			// Open file
			fh, err := os.Open(file)
			if err != nil {
				fmt.Println("error opening file")
				return err
			}
			defer fh.Close()

			//io.copy
			_, err = io.Copy(fileWriter, fh)
			if err != nil {
				return err
			}

			contentType := bodyWriter.FormDataContentType()
			log.Debug().Str("contentType", contentType).Msg("contentType")
			bodyWriter.Close()

			r := &Request{
				Type: "POST",
				Path: "chart",
				Body: bodyBuf.Bytes(),
			}

			serviceUrl := viper.GetString("serviceurl")
			apiVersion := api.APIVersion + "/"
			if serviceUrl == "" {
				serviceUrl = "localhost:8080/"
			}
			Url := "http://" + serviceUrl + "/" + apiVersion + r.Path
			body := bytes.NewReader(r.Body)
			log.Debug().Str("Type", r.Type).Str("url", Url).Msg("Request")
			request, err := http.NewRequest(r.Type, Url, body)
			if err != nil {
				return err
			}
			token, err := getToken()
			if err != nil {
				return err
			}
			Authorization := fmt.Sprintf("Bearer %s", token)

			request.Header.Add("Authorization", Authorization)
			request.Header.Add("Content-Type", contentType)

			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			log.Debug().Int("Code", resp.StatusCode).Bytes("Body", respBody).Msg("ok")
			fmt.Println("Upload file Success")
			return nil
		},
	}
}

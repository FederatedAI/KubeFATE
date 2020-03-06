package cli

import (
	"bytes"
	"encoding/json"
	"fate-cloud-agent/pkg/api"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
)

func VersionCommand() *cli.Command {
	return &cli.Command{
		Name: "version",
		Flags: []cli.Flag{
		},
		Usage: "show kubefate version",
	    Action: func(c *cli.Context) error {

			r := &request{
				Type: "GET",
				Path: "version",
				Body: nil,
			}

			serviceUrl := viper.GetString("serviceurl")
			apiVersion := api.ApiVersion + "/"
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
				return  err
			}
			Authorization := fmt.Sprintf("Bearer %s", token)

			request.Header.Add("Authorization", Authorization)

			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			type VersionResultMsg struct {
				Msg string
				Version  string
			}

			VersionResult := new(VersionResultMsg)

			err = json.Unmarshal(respBody, &VersionResult)
			if err != nil {
				return err
			}

			log.Debug().Int("Code", resp.StatusCode).Bytes("Body", respBody).Msg("ok")

			fmt.Printf("kubefate service version=%s\n", VersionResult.Version)
			fmt.Printf("kubefate commandLine version=%s\n", api.ServiceVersion)
			return nil
		},
	}
}

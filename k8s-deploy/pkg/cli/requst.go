/*
 * Copyright 2019-2020 VMware, Inc.
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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func getToken() (string, error) {

	serviceurl := viper.GetString("serviceurl")

	loginUrl := "http://" + serviceurl + "/v1/user/login"

	login := map[string]string{
		"username": viper.GetString("user.username"),
		"password": viper.GetString("user.password"),
	}

	loginJsonB, err := json.Marshal(login)

	body := bytes.NewReader(loginJsonB)
	Request, err := http.NewRequest("POST", loginUrl, body)
	if err != nil {
		return "", err
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(Request)
	if err != nil {
		return "", err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	Result := map[string]interface{}{}

	err = json.Unmarshal(rbody, &Result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprint(Result["message"]))
	}

	token := fmt.Sprint(Result["token"])

	log.Debug().Str("token", token).Msg("get token success")
	return token, nil
}

type Request struct {
	Type string
	Path string
	Body []byte
}

type Response struct {
	Code int
	Body []byte
}

func Send(r *Request) (*Response, error) {
	serviceUrl := viper.GetString("serviceurl")
	apiVersion := api.ApiVersion + "/"
	if serviceUrl == "" {
		serviceUrl = "localhost:8080/"
	}
	Url := "http://" + serviceUrl + "/" + apiVersion + r.Path
	body := bytes.NewReader(r.Body)
	log.Debug().Str("Type", r.Type).Str("url", Url).Str("Body", string(r.Body)).Msg("Request")
	request, err := http.NewRequest(r.Type, Url, body)
	if err != nil {
		return nil, err
	}
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	Authorization := fmt.Sprintf("Bearer %s", token)

	request.Header.Add("Authorization", Authorization)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Code: resp.StatusCode,
		Body: respBody,
	}, nil
}

type Result struct {
	Data []*modules.Job
	Msg  string
}

func (r *Response) Unmarshal() *Result {
	res := new(Result)
	_ = json.Unmarshal(r.Body, &res)
	return res
}

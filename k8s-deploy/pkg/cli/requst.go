package cli

import (
	"bytes"
	"encoding/json"
	"fate-cloud-agent/pkg/api"
	"fate-cloud-agent/pkg/db"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func getKubefate() {

	serviceurl := viper.GetString("serviceurl")

	//loginUrl := "http://" + serviceurl + "/v1"
	loginUrl := "http://" + serviceurl + "/v1/user/login"

	body := bytes.NewReader([]byte("{\"username\": \"admin\",\"password\": \"admin\"}"))
	request, err := http.NewRequest("POST", loginUrl, body)
	if err != nil {
		panic(err)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", rbody)
}

func getToken() string {

	serviceurl := viper.GetString("serviceurl")

	loginUrl := "http://" + serviceurl + "/v1/user/login"

	login := map[string]string{
		"username": viper.GetString("user.username"),
		"password": viper.GetString("user.password"),
	}

	loginJsonB, err := json.Marshal(login)

	body := bytes.NewReader(loginJsonB)
	request, err := http.NewRequest("POST", loginUrl, body)
	if err != nil {
		panic(err)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{}

	err = json.Unmarshal(rbody, &result)
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(result["token"])
}

type request struct {
	Type string
	Path string
	Body []byte
}

type Response struct {
	Code int
	Body []byte
}

func Send(r *request) (*Response, error) {
	serviceUrl := viper.GetString("serviceurl")
	apiVersion := api.ApiVersion+"/"
	if serviceUrl == "" {
		serviceUrl = "localhost:8080/"
	}
	Url := "http://" + serviceUrl +"/"+ apiVersion + r.Path
	body := bytes.NewReader(r.Body)
	log.Debug().Str("Type", r.Type).Str("url", Url).Str("Body",string(r.Body)).Msg("Request")
	request, err := http.NewRequest(r.Type, Url, body)
	if err != nil {
		return nil, err
	}

	Authorization := fmt.Sprintf("Bearer %s", getToken())

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

type result struct {
	Data []*db.Job
	Msg  string
}

func (r *Response) Unmarshal() *result {
	res := new(result)
	_ = json.Unmarshal(r.Body, &res)
	return res
}

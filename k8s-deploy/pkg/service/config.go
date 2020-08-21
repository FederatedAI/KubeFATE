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
package service

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"k8s.io/client-go/rest"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

func MapToConfig(m map[string]interface{}, templates string) (string, error) {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("fate-values-templates").Funcs(funcMap()).Option("missingkey=zero").Parse(string(templates)))

	// Execute the template for each recipient.

	var buf strings.Builder
	err := t.Execute(&buf, m)
	if err != nil {
		log.Error().Msg("executing template:" + err.Error())
		return "", err
	}
	s := strings.ReplaceAll(buf.String(), "<no value>", "")
	return s, nil

}

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	return f
}

func InitKubeConfig() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("~/.kube/config", []byte(config.String()), os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

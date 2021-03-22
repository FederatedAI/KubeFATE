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
package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func TestInitViper(t *testing.T) {
	InitViper()

	defaultConfig := []byte(`db:
  file: file::memory:?cache=shared
  type: sqlite
log:
  level: info
  nocolor: "false"
repo:
  name: kubefate
  url: https://federatedai.github.io/KubeFATE
server:
  address: 0.0.0.0
  port: "8080"
serviceurl: localhost:8080
user:
  password: admin
  username: admin
`)
	want := defaultConfig

	config := viper.AllSettings()
	got, err := yaml.Marshal(config)
	if (err != nil) != false {
		t.Errorf("GetPodList() error = %v, wantErr %v", err, false)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AllKeys() = \n%v \nwant \n%v", got, want)
	}

	fmt.Println(viper.GetString("db.type"))
}

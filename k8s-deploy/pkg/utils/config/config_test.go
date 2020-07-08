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
package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestConfig_DirExists(t *testing.T) {
	tmpDir := os.TempDir()
	exists := DirExists(tmpDir)
	if exists != true {
		t.Errorf("%s exists but DirExists return false \n", tmpDir)
	}

	// construct a random dir path
	tmpDir = time.Now().Format(time.RFC3339Nano)
	exists = DirExists(tmpDir)
	if exists != false {
		t.Errorf("%s does not exist but DirExists return true \n", tmpDir)
	}
}

func TestConfig_InitViper(t *testing.T) {
	_ = InitViper()

	viper.AddConfigPath("../../../")

	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("Fatal error config file: %s \n", err)
	}

	result := viper.Get("mongo")
	if result == "" {
		t.Errorf("Can not read mongo")
	}

	t.Log(result)
	result = viper.Get("mongo.url")
	t.Log(result)
}

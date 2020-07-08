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

package orm

import (
	"os"
	"reflect"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/spf13/viper"
)

func TestGetDbConfig(t *testing.T) {
	InitConfigForTest()
	tests := []struct {
		name string
		want *DbConfig
	}{
		// TODO: Add test cases.
		{
			name: "Environment variables as configuration files",
			want: &DbConfig{
				DbType:   os.Getenv("FATECLOUD_DB_TYPE"),
				Host:     os.Getenv("FATECLOUD_DB_HOST"),
				Port:     os.Getenv("FATECLOUD_DB_PORT"),
				Name:     os.Getenv("FATECLOUD_DB_NAME"),
				Username: os.Getenv("FATECLOUD_DB_USERNAME"),
				Password: os.Getenv("FATECLOUD_DB_PASSWORD"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDbConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDbConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
}

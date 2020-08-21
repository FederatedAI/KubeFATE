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
	"os"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"

	"github.com/spf13/viper"
)

func TestGetChartPath(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	_ = os.Setenv("FATECLOUD_CHART_PATH", "./")
	type args struct {
		name    string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				name:    "fate",
				version: "v1.2.0",
			},
			want: viper.GetString("chart.path") + "fate/v1.2.0/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChartPath(tt.args.name); got != tt.want {
				t.Errorf("GetChartPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

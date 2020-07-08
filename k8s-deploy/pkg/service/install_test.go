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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestInstall(t *testing.T) {

	_ = config.InitViper()
	viper.AddConfigPath("../../")
	_ = viper.ReadInConfig()
	logging.InitLog()
	err := RepoAddAndUpdate()
	if err != nil {
		panic(err)
	}
	type args struct {
		namespace    string
		name         string
		chartName    string
		chartVersion string
		value        *Value
	}
	tests := []struct {
		name    string
		args    args
		want    *releaseElement
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "install fate",
			args: args{
				namespace:    "fate-10000",
				name:         "fate-10000",
				chartName:    "fate",
				chartVersion: "v1.2.0",
				value:        &Value{Val: []byte(`{ "partyId":10000,"endpoint": { "ip":"192.168.100.123","port":30000}}`), T: "json"},
			},
			want: &releaseElement{
				Name:       "fate",
				Namespace:  "fate",
				Revision:   "1",
				Status:     "deployed",
				Chart:      "fate-1.2.0",
				AppVersion: "1.2.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Install(tt.args.namespace, tt.args.name, tt.args.chartName, tt.args.chartVersion, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Install() = %v, want %v", got, tt.want)
			}

		})
	}
}

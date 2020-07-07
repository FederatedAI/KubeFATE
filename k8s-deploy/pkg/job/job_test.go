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

// job
package job

import (
	"encoding/json"
	"fmt"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
	logging.InitLog()
}

func TestMsa(t *testing.T) {

	d := ClusterArgs{
		Name:         "fate-10000",
		Namespace:    "fate-10000",
		ChartName:    "fate",
		ChartVersion: "v1.3.0-a",
		Data:         []byte(`{"egg":{"count":3},"exchange":{"ip":"192.168.1.1","port":9370},"modules":["proxy","egg","fateboard","fateflow","federation","metaService","mysql","redis","roll","python"],"partyId":10000,"proxy":{"nodePort":30010,"type":"NodePort"}}`),
	}
	b, err := json.Marshal(d)
	if err != nil {
		log.Err(err).Msg("err")
	}

	fmt.Printf("%s", b)

}

func TestClusterInstall(t *testing.T) {
	InitConfigForTest()
	type args struct {
		clusterArgs *ClusterArgs
	}
	tests := []struct {
		name    string
		args    args
		want    *modules.Job
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "job install fate-8888",
			args: args{
				clusterArgs: &ClusterArgs{
					Name:         "fate-8888",
					Namespace:    "fate-8888",
					ChartName:    "fate",
					ChartVersion: "v1.3.0-a",
					Cover:        true,
					Data:         []byte(`{ "partyId":8888,"endpoint": { "ip":"192.168.100.123","port":30008}}`),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterInstall(tt.args.clusterArgs, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterInstall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterInstall() = %v, want %v", got, tt.want)
			}
			time.Sleep(60 * time.Second)
		})
	}
}

func TestClusterDelete(t *testing.T) {
	InitConfigForTest()
	type args struct {
		clusterId string
	}
	tests := []struct {
		name string
		args args
		want *modules.Job
	}{
		// TODO: Add test cases.
		{
			name: "delete",
			args: args{
				clusterId: "5029628c-8886-4907-bced-6dbe3553c7ef",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ClusterDelete(tt.args.clusterId, "test"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterDelete() = %v, want %v", got, tt.want)
			}
			time.Sleep(30 * time.Second)
		})
	}
}

func TestClusterUpdate(t *testing.T) {
	InitConfigForTest()
	type args struct {
		clusterArgs *ClusterArgs
		creator     string
	}
	tests := []struct {
		name    string
		args    args
		want    *modules.Job
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				clusterArgs: &ClusterArgs{
					Name:         "fate-9999",
					Namespace:    "fate-9999",
					ChartName:    "fate",
					ChartVersion: "v1.3.0-b",
					Cover:        false,
					Data:         []byte(`{"chartName":"fate","chartVersion":"v1.3.0-b","egg":{"count":1},"modules":["proxy","egg","federation","metaService","mysql","redis","roll","python"],"name":"fate-10000","namespace":"fate-10000","partyId":10000,"proxy":{"exchange":{"ip":"192.168.100.123","port":30000},"nodePort":30010,"type":"NodePort"}}`),
				},
				creator: "admin",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterUpdate(tt.args.clusterArgs, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterUpdate() = %v, want %v", got, tt.want)
			}
			time.Sleep(30 * time.Second)
		})
	}
}

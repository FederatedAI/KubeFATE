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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

func TestAbc(t *testing.T) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.String())
}

func TestGetPods(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		time.Sleep(time.Second)
		t.Run(tt.name, func(t *testing.T) {
			labelSelector := fmt.Sprintf("name=%s", tt.args.namespace)
			got, err := GetPods(tt.args.namespace, labelSelector)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("Namespace, Name, Status")
			for _, v := range got.Items {
				for _, vv := range v.Status.ContainerStatuses {
					fmt.Println(vv.State.String())
				}
				fmt.Printf("%s, %s, %s\n", v.Namespace, v.Name, v.Status.Phase)
			}
		})
	}
}

func TestGetPodStatus(t *testing.T) {
	InitConfigForTest()
	pods, _ := GetPods("fate-9999", "name=fate-9999")
	type args struct {
		pods *v1.PodList
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				pods: pods,
			},
			want: map[string]string{
				"client":         "Running",
				"clustermanager": "Running",
				"fateboard":      "Running",
				"mysql":          "Running",
				"nodemanager":    "Running",
				"python":         "Running",
				"rollsite":       "Running",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPodStatus(tt.args.pods); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
func InitConfigForTest() {
	_ = config.InitViper()
	viper.AddConfigPath("../../")
	_ = viper.ReadInConfig()
	logging.InitLog()
}

func TestCheckClusterStatus(t *testing.T) {
	type args struct {
		ClusterStatus map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "false",
			args: args{
				ClusterStatus: map[string]string{},
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				ClusterStatus: map[string]string{"python": "Running", "client": "Running", "fateboard": "Pending"},
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				ClusterStatus: map[string]string{"python": "Running", "client": "Running", "fateboard": "Running"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckClusterStatus(tt.args.ClusterStatus); got != tt.want {
				t.Errorf("CheckClusterStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

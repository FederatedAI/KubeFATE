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
	"testing"
	"time"

	"k8s.io/client-go/rest"
)

func TestAbc(t *testing.T) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.String())
}

func Test_checkClusterStatus(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "cluster is running",
			args: args{
				namespace: "fate-10000",
				name:      "fate-10000",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckClusterStatus(tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkClusterStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkClusterStatus() = %v, want %v", got, tt.want)
			}
		})
	}
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

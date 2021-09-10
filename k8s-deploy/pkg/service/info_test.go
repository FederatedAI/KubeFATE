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

package service

import "testing"

func TestCheckClusterInfoStatus(t *testing.T) {
	type args struct {
		ClusterInfoStatus map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "nil",
			args: args{
				ClusterInfoStatus: map[string]interface{}{},
			},
			want: false,
		},
		{
			name: "status-nil",
			args: args{
				ClusterInfoStatus: map[string]interface{}{
					"status": map[string]interface{}{},
				},
			},
			want: false,
		},
		{
			name: "deployments-nil",
			args: args{
				ClusterInfoStatus: map[string]interface{}{
					"status": map[string]interface{}{
						"deployments": map[string]string{},
					},
				},
			},
			want: false,
		},
		{
			name: "python-false",
			args: args{
				ClusterInfoStatus: map[string]interface{}{
					"status": map[string]interface{}{
						"deployments": map[string]string{
							"client":         "Available",
							"clustermanager": "Available",
							"mysql":          "Available",
							"python":         "Progressing",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "all-Available",
			args: args{
				ClusterInfoStatus: map[string]interface{}{
					"status": map[string]interface{}{
						"deployments": map[string]string{
							"client":         "Available",
							"clustermanager": "Available",
							"mysql":          "Available",
							"python":         "Available",
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckClusterInfoStatus(tt.args.ClusterInfoStatus); got != tt.want {
				t.Errorf("CheckClusterInfoStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

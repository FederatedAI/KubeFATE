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
package cli

import (
	"encoding/json"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
)

func TestJob_outPutInfo(t *testing.T) {

	type args struct {
		result interface{}
	}
	tests := []struct {
		name    string
		c       *Job
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.

		{name: "",
			c: &Job{},
			args: args{
				result: func() *JobResult {
					job := modules.Job{}
					_ = json.Unmarshal([]byte(`{"CreatedAt":"2020-07-17T13:40:53+08:00","DeletedAt":null,"ID":19,"UpdatedAt":"2020-07-17T13:40:54+08:00","cluster_id":"8c1cec52-7afc-4716-9b55-158f66c71b34","creator":"test","end_time":"2020-07-17T13:40:54+08:00","method":"ClusterInstall","result":"template: fate-values-templates:18:20: executing \"fate-values-templates\" at \u003c.istio.enabled\u003e: nil pointer evaluating interface {}.enabled","start_time":"2020-07-17T13:40:53+08:00","status":"Failed","sub_jobs":null,"time_limit":3600000000000,"uuid":"6d1c3a99-3e91-4ae7-baf8-5b2864640d2d"}`), &job)
					return &JobResult{Data: &job, Msg: ""}
				}(),
			}, wantErr: false},
		{name: "",
			c: &Job{},
			args: args{
				result: func() *JobResult {
					job := modules.Job{}
					_ = json.Unmarshal([]byte(`{"CreatedAt":"2020-07-20T14:55:56+08:00","DeletedAt":null,"ID":24,"UpdatedAt":"2020-07-20T14:56:22+08:00","cluster_id":"cac7cf67-6b55-4619-a3fd-5bb8aef515f1","creator":"test","end_time":"2020-07-20T14:56:22+08:00","method":"ClusterInstall","result":"Cluster install success","start_time":"2020-07-20T14:55:55+08:00","status":"Success","sub_jobs":{"client":{"EndTime":"2020-07-20T14:55:56.9104784+08:00","ModuleName":"client","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"},"clustermanager":{"EndTime":"2020-07-20T14:56:01.9773077+08:00","ModuleName":"clustermanager","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"},"fateboard":{"EndTime":"2020-07-20T14:55:56.9104784+08:00","ModuleName":"fateboard","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"},"mysql":{"EndTime":"2020-07-20T14:56:01.9773077+08:00","ModuleName":"mysql","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"},"nodemanager":{"EndTime":"2020-07-20T14:56:01.9773077+08:00","ModuleName":"nodemanager","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"},"python":{"EndTime":"2020-07-20T14:55:56.9104784+08:00","ModuleName":"python","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"},"rollsite":{"EndTime":"2020-07-20T14:56:01.9773077+08:00","ModuleName":"rollsite","ModulesPodStatus":"Running","StartTime":"2020-07-20T14:55:55.969774+08:00","Status":"Success"}},"time_limit":3600000000000,"uuid":"188e73a9-8d1d-4b4e-9b32-e09bfcbb1d9a"}`), &job)
					return &JobResult{Data: &job, Msg: ""}
				}(),
			}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{}
			if err := c.outPutInfo(tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("Job.outPutInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
	"reflect"
	"testing"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var CLUSTERID string

func InitConfigForTest() {
	_ = config.InitViper()
	viper.AddConfigPath("../../")
	_ = viper.ReadInConfig()
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
	log.Level(zerolog.DebugLevel)
	mysql := new(orm.Mysql)
	_ = mysql.Setup()
	modules.DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	new(modules.User).InitTable()
	new(modules.Cluster).InitTable()
	new(modules.HelmChart).InitTable()
	new(modules.Job).InitTable()

	type args struct {
		clusterArgs *ClusterArgs
	}
	tests := []struct {
		name    string
		args    args
		want    modules.JobStatus
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "job install error and rollback",
			args: args{
				clusterArgs: &ClusterArgs{
					Name:         "fate-9999",
					Namespace:    "fate-9999",
					ChartName:    "fate",
					ChartVersion: "v1.4.1",
					Cover:        false,
					Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.1","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":2,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"192.168.10.1","partyPort":30010}],"type":"NodePort"}}`),
				},
			},
			want:    modules.JobStatusFailed,
			wantErr: false,
		},
		{
			name: "job install fate-9999 success",

			args: args{
				clusterArgs: &ClusterArgs{
					Name:         "fate-9999",
					Namespace:    "fate-9999",
					ChartName:    "fate",
					ChartVersion: "v1.4.0",
					Cover:        true,
					Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.0","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":0,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"10.117.32.179","partyPort":30010}],"type":"NodePort"}}`),
				},
			},
			want:    modules.JobStatusSuccess,
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
			if got == nil {
				t.Errorf("ClusterUpdate() = %v, want %s", got, "not nil")
				return
			}
			CLUSTERID = got.Uuid
			time.Sleep(5 * time.Second)
			for got.Status == modules.JobStatusRunning || got.Status == modules.JobStatusPending {
				time.Sleep(5 * time.Second)
			}
			if got.Status != tt.want {
				t.Errorf("ClusterInstall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterUpdate(t *testing.T) {
	InitConfigForTest()
	log.Level(zerolog.DebugLevel)
	mysql := new(orm.Mysql)
	_ = mysql.Setup()
	modules.DB = orm.DBCLIENT
	//DB.LogMode(true)
	type args struct {
		clusterArgs *ClusterArgs
		creator     string
	}
	tests := []struct {
		name    string
		args    args
		want    modules.JobStatus
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
					ChartVersion: "v1.4.0",
					Cover:        false,
					Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.0","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":2,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"192.168.10.1","partyPort":30010}],"type":"NodePort"}}`),
				},
				creator: "admin",
			},
			want:    modules.JobStatusSuccess,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				clusterArgs: &ClusterArgs{
					Name:         "fate-9999",
					Namespace:    "fate-9999",
					ChartName:    "fate",
					ChartVersion: "v1.4.0",
					Cover:        false,
					Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.0","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":2,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"192.168.10.1","partyPort":30010}],"type":"NodePort"}}`),
				},
				creator: "admin",
			},
			want:    modules.JobStatusSuccess,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				clusterArgs: &ClusterArgs{
					Name:         "fate-9999",
					Namespace:    "fate-9999",
					ChartName:    "fate",
					ChartVersion: "v1.4.1-a",
					Cover:        false,
					Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.1-a","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":2,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"192.168.10.1","partyPort":30010}],"type":"NodePort"}}`),
				},
				creator: "admin",
			},
			want:    modules.JobStatusSuccess,
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
			if got == nil {
				return
			}
			time.Sleep(5 * time.Second)
			for got.Status == modules.JobStatusRunning || got.Status == modules.JobStatusPending {
				time.Sleep(5 * time.Second)
			}
			if got.Status != tt.want {
				t.Errorf("ClusterInstall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterDelete(t *testing.T) {
	InitConfigForTest()
	log.Level(zerolog.DebugLevel)
	mysql := new(orm.Mysql)
	_ = mysql.Setup()
	modules.DB = orm.DBCLIENT
	//DB.LogMode(true)
	type args struct {
		clusterId string
		creator   string
	}
	tests := []struct {
		name    string
		args    args
		want    modules.JobStatus
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete",
			args: args{
				clusterId: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				clusterId: CLUSTERID,
				creator:   "test",
			},
			want:    modules.JobStatusSuccess,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterDelete(tt.args.clusterId, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("ClusterDelete() = %v, want %s", got, "not nil")
				return
			}
			time.Sleep(5 * time.Second)
			for got.Status == modules.JobStatusRunning || got.Status == modules.JobStatusPending {
				time.Sleep(5 * time.Second)
			}
			if got.Status != tt.want {
				t.Errorf("ClusterInstall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster(t *testing.T) {
	InitConfigForTest()
	log.Level(zerolog.DebugLevel)
	mysql := new(orm.Mysql)
	_ = mysql.Setup()
	modules.DB = orm.DBCLIENT
	//DB.LogMode(true)
	// DropTable Table
	new(modules.User).DropTable()
	new(modules.Cluster).DropTable()
	new(modules.HelmChart).DropTable()
	new(modules.Job).DropTable()
	// Create Table
	new(modules.User).InitTable()
	new(modules.Cluster).InitTable()
	new(modules.HelmChart).InitTable()
	new(modules.Job).InitTable()

	clusterArgs := &ClusterArgs{
		Name:         "fate-9999",
		Namespace:    "fate-9999",
		ChartName:    "fate",
		ChartVersion: "v1.4.0",
		Cover:        true,
		Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.0","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":0,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"10.117.32.179","partyPort":30010}],"type":"NodePort"}}`),
	}

	var clusterId string

	t.Run("ClusterInstall", func(t *testing.T) {
		got, err := ClusterInstall(clusterArgs, "test")
		if (err != nil) != false {
			t.Errorf("ClusterInstall() error = %v, wantErr %v", err, false)
			return
		}
		if got == nil {
			t.Errorf("ClusterInstall() = %v, want %s", got, "not nil")
			return
		}
		clusterId = got.ClusterId
		for got.Status == modules.JobStatusRunning || got.Status == modules.JobStatusPending {
			time.Sleep(5 * time.Second)
		}
		if got.Status != modules.JobStatusSuccess {
			t.Errorf("ClusterInstall() = %v, want %v", got, modules.JobStatusSuccess)
		}
	})

	clusterArgsUpdate := &ClusterArgs{
		Name:         "fate-9999",
		Namespace:    "fate-9999",
		ChartName:    "fate",
		ChartVersion: "v1.4.1-a",
		Cover:        false,
		Data:         []byte(`{"chartName":"fate","chartVersion":"v1.4.1-a","modules":["rollsite","clustermanager","nodemanager","mysql","python","client"],"mysql":{"accessMode":"ReadWriteOnce","database":"eggroll_meta","existingClaim":"","ip":"mysql","nodeSelector":{},"password":"fate_dev","port":3306,"size":"1Gi","storageClass":"mysql","subPath":"","user":"fate"},"name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":2,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false,"pullPolicy":null,"python":{"fateflowNodePort":30109,"fateflowType":"NodePort","nodeSelector":{}},"registry":"","rollsite":{"exchange":{"ip":"192.168.1.1","port":30000},"nodePort":30009,"nodeSelector":{},"partyList":[{"partyId":10000,"partyIp":"192.168.10.1","partyPort":30010}],"type":"NodePort"}}`),
	}

	t.Run("ClusterUpdate", func(t *testing.T) {
		got, err := ClusterUpdate(clusterArgsUpdate, "test")
		if (err != nil) != false {
			t.Errorf("ClusterUpdate() error = %v, wantErr %v", err, false)
			return
		}
		if got == nil {
			t.Errorf("ClusterUpdate() = %v, want %s", got, "not nil")
			return
		}
		time.Sleep(5 * time.Second)
		for got.Status == modules.JobStatusRunning || got.Status == modules.JobStatusPending {
			time.Sleep(5 * time.Second)
		}
		if got.Status != modules.JobStatusSuccess {
			t.Errorf("ClusterUpdate() = %v, want %v", got, modules.JobStatusSuccess)
		}
	})

	t.Run("ClusterDelete", func(t *testing.T) {
		got, err := ClusterDelete(clusterId, "test")
		if (err != nil) != false {
			t.Errorf("ClusterDelete() error = %v, wantErr %v", err, false)
			return
		}
		if got == nil {
			t.Errorf("ClusterDelete() = %v, want %s", got, "not nil")
			return
		}
		time.Sleep(5 * time.Second)
		for got.Status == modules.JobStatusRunning || got.Status == modules.JobStatusPending {
			time.Sleep(5 * time.Second)
		}
		if got.Status != modules.JobStatusSuccess {
			t.Errorf("ClusterDelete() = %v, want %v", got, modules.JobStatusSuccess)
		}
	})
}

func Test_generateSubJobs(t *testing.T) {
	var job modules.Job
	type args struct {
		job           *modules.Job
		ClusterStatus map[string]string
	}
	tests := []struct {
		name string
		args args
		want modules.SubJobs
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				job:           &job,
				ClusterStatus: map[string]string{"client": "Waiting", "clustermanager": "Waiting", "fateboard": "Waiting", "mysql": "Waiting", "nodemanager": "Waiting", "python": "Waiting", "rollsite": "Waiting"},
			},
			want: nil,
		},
		{
			name: "",
			args: args{
				job:           &job,
				ClusterStatus: map[string]string{"client": "Running", "clustermanager": "Running", "fateboard": "Running", "mysql": "Running", "nodemanager": "Running", "python": "Running", "rollsite": "Running"},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateSubJobs(tt.args.job, tt.args.ClusterStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateSubJobs() = %v, want %v", got, tt.want)
			}
			time.Sleep(time.Second * 5)
		})
	}
}

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

package modules

import (
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
)

func TestGetFateChart(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)
	// Create Table
	e := &HelmChart{}

	e.InitTable()
	// Drop Table
	defer e.DropTable()
	type args struct {
		chartName    string
		chartVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    *HelmChart
		wantErr bool
	}{
		{
			name: "not is existed",
			args: args{
				chartName:    "fate",
				chartVersion: "v1.4.0",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "is existed",
			args: args{
				chartName:    "fate",
				chartVersion: "v1.4.0",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFateChart(tt.args.chartName, tt.args.chartVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFateChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("GetFateChart() = %v, want %v", got, tt.want)
			}
		})
	}
}

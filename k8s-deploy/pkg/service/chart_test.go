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
	"reflect"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"

	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/chart"
)

func TestFateChart_save(t *testing.T) {
	type fields struct {
		version   string
		Chart     *chart.Chart
		HelmChart *modules.HelmChart
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FateChart{
				HelmChart: tt.fields.HelmChart,
			}
			if err := fc.save(); (err != nil) != tt.wantErr {
				t.Errorf("FateChart.save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

func TestGetFateChart(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	_ = os.Setenv("FATECLOUD_CHART_PATH", "../../")
	type args struct {
		name    string
		version string
	}
	var tests = []struct {
		name    string
		args    args
		want    *FateChart
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				name:    "fate",
				version: "v1.2.0",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFateChart(tt.args.name, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFateChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFateChart() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFateChart_read(t *testing.T) {
	InitConfigForTest()
	type fields struct {
		HelmChart *modules.HelmChart
	}
	type args struct {
		name    string
		version string
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			fields: fields{
				HelmChart: new(modules.HelmChart),
			},
			args: args{
				name:    "fate",
				version: "v1.2.0",
			},
			want:    "v1.2.0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FateChart{
				HelmChart: tt.fields.HelmChart,
			}
			got, err := fc.read(tt.args.name, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("FateChart.read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Version != tt.want {
				t.Errorf("FateChart.read() Version = %v, want %v", got.Version, tt.want)
			}
		})
	}
}

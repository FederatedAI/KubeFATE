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
package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindHelmCharts(t *testing.T) {
	InitConfigForTest()
	job := &HelmChart{}
	results, _ := Find(job)
	t.Log(ToJson(results))
}
func TestHelmChart_FindHelmByVersion(t *testing.T) {
	InitConfigForTest()
	type args struct {
		version string
		name    string
	}
	tests := []struct {
		name    string
		args    args
		want    *HelmChart
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "read",
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
			got, err := FindHelmByNameAndVersion(tt.args.name, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("HelmChart.FindHelmByNameAndVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HelmChart.FindHelmByNameAndVersion() = %+v, want %v", got, tt.want)
			}
		})
	}
}

func TestChartDeleteAll(t *testing.T) {
	InitConfigForTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Error().Err(err).Msg("ConnectDb")
	}
	collection := db.Collection(new(HelmChart).getCollection())
	filter := bson.D{}
	r, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("DeleteMany")
	}
	if r.DeletedCount == 0 {
		log.Error().Msg("this record may not exist(DeletedCount==0)")
	}
	fmt.Println(r)
	return
}

func TestFindHelmChartList(t *testing.T) {
	InitConfigForTest()
	tests := []struct {
		name    string
		want    []*HelmChart
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindHelmChartList()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindHelmChartList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range got {
				t.Logf("%+v\n", v)
			}
		})
	}
}

func TestChartSave(t *testing.T) {
	InitConfigForTest()
	type args struct {
		helmChart *HelmChart
	}
	var tests = []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				helmChart: &HelmChart{
					Uuid:           "4444",
					Name:           "fate444",
					Chart:          "fate",
					Values:         "fate",
					ValuesTemplate: "ddd",
					Templates:      nil,
					Version:        "v1.2.4",
				},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChartSave(tt.args.helmChart)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChartSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChartSave() = %v, want %v", got, tt.want)
			}
		})
	}
}

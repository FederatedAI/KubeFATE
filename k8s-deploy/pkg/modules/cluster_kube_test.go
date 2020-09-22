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
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestCluster_GetClusterStatus(t *testing.T) {
	type fields struct {
		Uuid         string
		Name         string
		NameSpace    string
		ChartName    string
		ChartVersion string
		Values       string
		Spec         MapStringInterface
		Revision     int8
		HelmRevision int8
		ChartValues  MapStringInterface
		Status       ClusterStatus
		Info         MapStringInterface
		Model        gorm.Model
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				NameSpace: "fate-9999",
				Name:      "fate-9999",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Values:       tt.fields.Values,
				Spec:         tt.fields.Spec,
				Revision:     tt.fields.Revision,
				HelmRevision: tt.fields.HelmRevision,
				ChartValues:  tt.fields.ChartValues,
				Status:       tt.fields.Status,
				Info:         tt.fields.Info,
				Model:        tt.fields.Model,
			}
			got, err := e.GetClusterStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cluster.GetClusterStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cluster.GetClusterStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

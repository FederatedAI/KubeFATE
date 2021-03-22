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

package orm

import (
	"reflect"
	"testing"
)

func Test_getDbType(t *testing.T) {
	type args struct {
		Type string
	}
	tests := []struct {
		name    string
		args    args
		want    Database
		wantErr bool
	}{
		{
			name: "mysql",
			args: args{
				Type: "mysql",
			},
			want:    &Mysql{},
			wantErr: false,
		},
		{
			name: "sqlite",
			args: args{
				Type: "sqlite",
			},
			want:    &Sqlite{},
			wantErr: false,
		},
		{
			name: "unknown",
			args: args{
				Type: "unknown",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDbType(tt.args.Type)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDbType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDbType() = %v, want %v", got, tt.want)
			}
		})
	}
}

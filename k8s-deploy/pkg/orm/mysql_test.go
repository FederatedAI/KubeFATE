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

package orm

import "testing"

func TestMysql_Setup(t *testing.T) {
	InitConfigForTest()
	tests := []struct {
		name    string
		e       *Mysql
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "The test database can be connected normally",
			e:       new(Mysql),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Mysql{}
			if err := e.Setup(); (err != nil) != tt.wantErr {
				t.Errorf("Mysql.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

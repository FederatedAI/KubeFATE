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

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
)

func TestUserStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    UserStatus
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Deprecate",
			s:    Deprecate_u,
			want: "Deprecate",
		},
		{
			name: "Available",
			s:    Available_u,
			want: "Available",
		},
		{
			name: "",
			s:    0,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("UserStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       UserStatus
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Deprecate",
			s:       Deprecate_u,
			want:    []byte("\"Deprecate\""),
			wantErr: false,
		},
		{
			name:    "Available",
			s:       Available_u,
			want:    []byte("\"Available\""),
			wantErr: false,
		},
		{
			name:    "",
			s:       0,
			want:    []byte("\"\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStatus.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encryption(t *testing.T) {
	type args struct {
		plaintext string
		salt      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				plaintext: "123",
				salt:      "456",
			},
			want: "9a5f5c39f26968a2b8c359eb03f5b5bcb38123c921cf2e078c31c8563bf52f900a1e3076270579d932c4db02fac68052ce52bae57f513afb24f9eb200b0efa72c74ec829083ba1fdd6830539741de4b7a4e04e6cef0f787664d212dfeca8af52335f71906fba07ed1009518634587e7c5be6fbb0cdd4142edbfccb37119e4831114b5198a4eef6a2874d6bef89f97652d9111e4162e69dc540609e063b660e718511040679f7bff0ab44143d7b85b690eea88309ee57b89c3094b2db55f7226e333d065c712c8576eec3960d424e5b8395d1753d63c60fdbf2200b4a0fe0007073b907affa73e5b57abc1e1a965487618f27e97d1387c431c0c66352e7955f68",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encryption(tt.args.plaintext, tt.args.salt); got != tt.want {
				t.Errorf("encryption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsValid(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &User{}
	e.InitTable()
	// Drop Table
	defer e.DropTable()

	//Insert
	e = NewUser("admin", "admin", "admin@admin.admin")
	Id, err := e.Insert()
	if err != nil {
		t.Errorf("User.Insert() error = %v", err)
		return
	}
	if Id != 1 {
		t.Errorf("User.Insert() got = %d, want = %d", Id, 1)
		return
	}

	e = &User{
		Username: "admin",
		Password: "admin",
	}
	if ok := e.IsValid(); !ok {
		t.Errorf("User.IsValid() got = %v, want = %v", ok, true)
		return
	}

	e = &User{
		Username: "admin",
		Password: "admin-error",
	}
	if ok := e.IsValid(); ok {
		t.Errorf("User.IsValid() got = %v, want = %v", ok, false)
		return
	}

}

func TestUser_IsExisted(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &User{}
	e.InitTable()
	// Drop Table
	defer e.DropTable()

	//Insert
	e = NewUser("admin", "admin", "admin@admin.admin")
	Id, err := e.Insert()
	if err != nil {
		t.Errorf("User.Insert() error = %v", err)
		return
	}
	if Id != 1 {
		t.Errorf("User.Insert() got = %d, want = %d", Id, 1)
		return
	}

	e = &User{
		Username: "admin",
		Password: "admin",
	}
	if ok := e.IsExisted(); !ok {
		t.Errorf("User.IsExisted() got = %v, want = %v", ok, true)
		return
	}

}

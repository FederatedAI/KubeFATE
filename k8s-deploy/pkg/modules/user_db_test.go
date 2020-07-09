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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func TestUser(t *testing.T) {
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

	want := e
	e = &User{Uuid: want.Uuid}
	got, err := e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err)
		return
	}
	if got.Username != want.Username || e.IsValid() || got.Uuid != want.Uuid || got.Email != want.Email || got.Password != "" {
		t.Errorf("User.Get() where Username=admin  got = %v, wat = %v", got, want)
		return
	}
	e = &User{Username: "admin"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err)
		return
	}
	if got.Username != want.Username || e.IsValid() || got.Uuid != want.Uuid || got.Email != want.Email || got.Password != "" {
		t.Errorf("User.Get() where Username=admin  got = %v, wat = %v", got, want)
		return
	}
	e = &User{Email: want.Email}
	got, err = e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err)
		return
	}
	if got.Username != want.Username || e.IsValid() || got.Uuid != want.Uuid || got.Email != want.Email || got.Password != "" {
		t.Errorf("User.Get() where Username=admin  got = %v, wat = %v", got, want)
		return
	}
	e = &User{Status: want.Status}
	got, err = e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err)
		return
	}
	if got.Username != want.Username || e.IsValid() || got.Uuid != want.Uuid || got.Email != want.Email || got.Password != "" {
		t.Errorf("User.Get() where Username=admin  got = %v, wat = %v", got, want)
		return
	}

	// Insert
	e = NewUser("root", "root", "root@root.root")
	Id, err = e.Insert()
	if err != nil {
		t.Errorf("User.Insert() error = %v", err)
		return
	}
	if Id != 2 {
		t.Errorf("User.Insert() got = %d, want = %d", Id, 2)
		return
	}

	// Insert
	e = NewUser("guest", "", "guest@guest.guest")
	Id, err = e.Insert()
	if err != nil {
		t.Errorf("User.Insert() error = %v", err)
		return
	}
	if Id != 3 {
		t.Errorf("User.Insert() got = %d, want = %d", Id, 3)
		return
	}

	e = &User{}
	gots, err := e.GetList()
	if err != nil {
		t.Errorf("User.GetList() error = %v", err)
		return
	}
	if len(gots) != 3 {
		t.Errorf("User.GetList() len(got) = %d, want = %d", len(gots), 3)
		return
	}
	e = &got
	gots, err = e.GetList()
	if err != nil {
		t.Errorf("User.GetList() error = %v", err)
		return
	}
	if len(gots) != 1 {
		t.Errorf("User.GetList() len(got) = %d, want = %d", len(gots), 1)
		return
	}

	// Update
	e = NewUser("admin-update", "admin-update", "admin-update@admin.admin")
	got, err = e.Update(1)
	if err != nil {
		t.Errorf("User.Update() error = %v", err)
		return
	}
	if got.Username != e.Username || got.Email != e.Email || e.IsValid() || got.Uuid != e.Uuid || got.Password != "" {
		t.Errorf("User.Update() got = %+v, want = %+v", got, e)
		return
	}

	want = e
	e = &User{}
	e.ID = 1
	got, err = e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err)
		return
	}
	if got.Username != want.Username || want.IsValid() || got.Uuid != want.Uuid || got.Email != want.Email || got.Password != "" {
		t.Errorf("User.Get() got = %v, wat = %v", got, want)
		return
	}

	e = &User{}
	success, err := e.DeleteById(1)
	if err != nil {
		t.Errorf("User.Delete() error = %v", err)
		return
	}
	if !success {
		t.Errorf("User.Delete() success = %v, wat = %v", success, true)
		return
	}
	e = &User{}
	success, err = e.DeleteById(0)
	if err != nil {
		t.Errorf("User.Delete() error = %v", err)
		return
	}

	e = &User{}
	gots, err = e.GetList()
	if err != nil {
		t.Errorf("User.GetList() error = %v", err)
		return
	}
	if len(gots) != 2 {
		t.Errorf("User.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}
}

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
}

func TestUser_Delete(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	type fields struct {
		Uuid     string
		Username string
		Password string
		Salt     string
		Email    string
		Status   UserStatus
		Model    gorm.Model
	}
	type args struct {
		id uint
	}
	var (
		tests = []struct {
			name        string
			fields      fields
			args        args
			wantSuccess bool
			wantErr     bool
		}{
			// TODO: Add test cases.
			{
				name:   "",
				fields: fields{},
				args: args{
					id: 0,
				},
				wantSuccess: false,
				wantErr:     true,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &User{
				Uuid:     tt.fields.Uuid,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				Salt:     tt.fields.Salt,
				Email:    tt.fields.Email,
				Status:   tt.fields.Status,
				Model:    tt.fields.Model,
			}
			gotSuccess, err := e.DeleteById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("User.Delete() = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

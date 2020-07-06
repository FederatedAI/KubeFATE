/*
 *  Copyright 2019-2020 VMware, Inc.
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
 */

package modules

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/spf13/viper"
	"testing"
)

func TestUser(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	db = orm.DBCLIENT
	//db.LogMode(true)

	// Create Table
	e := &User{}
	e.InitTable()
	// Drop Table
	defer e.DropTable()

	//Insert
	e = NewUser("admin", "admin", "admin@admin.admin")
	Id, err := e.Insert()
	if err != nil {
		t.Errorf("User.Insert() error = %v", err, )
		return
	}
	if Id != 1 {
		t.Errorf("User.Insert() got = %d, want = %d", Id, 1)
		return
	}

	want := e
	e = &User{Username: "admin"}
	got, err := e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err)
		return
	}
	if got.Username != want.Username || got.Password != want.Password || got.ID != want.ID || got.Email != want.Email {
		t.Errorf("User.Get() where Username=admin  got = %v, wat = %v", got, want)
		return
	}

	// Insert
	e = NewUser("root", "root", "root@root.root")
	Id, err = e.Insert()
	if err != nil {
		t.Errorf("User.Insert() error = %v", err, )
		return
	}
	if Id != 2 {
		t.Errorf("User.Insert() got = %d, want = %d", Id, 2)
		return
	}

	e = &User{}
	gots, err := e.GetList()
	if err != nil {
		t.Errorf("User.GetList() error = %v", err, )
		return
	}
	if len(gots) != 2 {
		t.Errorf("User.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}

	// Update
	e = NewUser("admin-update", "admin-update", "admin-update@admin.admin")
	got, err = e.Update(1)
	if err != nil {
		t.Errorf("User.Update() error = %v", err, )
		return
	}
	if got.Username != e.Username || got.Email != e.Email {
		t.Errorf("User.Update() got = %v, want = %v", got, e)
		return
	}

	e = &User{}
	e.ID = 1
	want = &got
	got, err = e.Get()
	if err != nil {
		t.Errorf("User.Get() error = %v", err, )
		return
	}
	if got.Username != want.Username || got.Password != want.Password || got.ID != want.ID || got.Email != want.Email {
		t.Errorf("User.Get() got = %v, wat = %v", got, want)
		return
	}

	e = &User{}
	success, err := e.Delete(1)
	if err != nil {
		t.Errorf("User.Delete() error = %v", err, )
		return
	}
	if !success {
		t.Errorf("User.Delete() success = %v, wat = %v", success, true)
		return
	}

	e = &User{}
	gots, err = e.GetList()
	if err != nil {
		t.Errorf("User.GetList() error = %v", err, )
		return
	}
	if len(gots) != 1 {
		t.Errorf("User.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}
}

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
}

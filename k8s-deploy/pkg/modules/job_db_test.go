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
	"fmt"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
	uuid "github.com/satori/go.uuid"
)

func TestJob(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Job{}

	e.InitTable()
	// Drop Table
	defer e.DropTable()

	//Insert
	e = NewJob("INSTALL", "admin", "")
	e.ClusterId = uuid.NewV4().String()
	Id, err := e.Insert()
	if err != nil {
		t.Errorf("Job.Insert() error = %v", err)
		return
	}
	if Id != 1 {
		t.Errorf("Job.Insert() got = %d, want = %d", Id, 1)
		return
	}
	// repeat insert
	Id, err = e.Insert()
	if err == nil {
		t.Error("Job.Insert() error = repeat insert")
		return
	}

	want := e
	e = &Job{Uuid: e.Uuid}
	got, err := e.Get()
	if err != nil {
		t.Errorf("Job.Get() error = %v", err)
		return
	}
	if got.Uuid != want.Uuid || got.Status != want.Status ||
		got.Method != want.Method || got.Creator != want.Creator || got.ClusterId != want.ClusterId {
		t.Errorf("Job.Get() where Name=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Job{Method: "INSTALL"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Job.Get() error = %v", err)
		return
	}
	if got.Uuid != want.Uuid || got.Status != want.Status ||
		got.Method != want.Method || got.Creator != want.Creator || got.ClusterId != want.ClusterId {
		t.Errorf("Job.Get() where Name=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Job{Creator: "admin"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Job.Get() error = %v", err)
		return
	}
	if got.Uuid != want.Uuid || got.Status != want.Status ||
		got.Method != want.Method || got.Creator != want.Creator || got.ClusterId != want.ClusterId {
		t.Errorf("Job.Get() where NameSpace=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Job{ClusterId: e.ClusterId}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Job.Get() error = %v", err)
		return
	}
	if got.Uuid != want.Uuid || got.Status != want.Status ||
		got.Method != want.Method || got.Creator != want.Creator || got.ClusterId != want.ClusterId {
		t.Errorf("Job.Get() where NameSpace=fate-9999  got = %v, wat = %v", got, want)
		return
	}

	// Insert
	e = NewJob("UPGRADE", "test-0", "")
	Id, err = e.Insert()
	if err != nil {
		t.Errorf("Job.Insert() error = %v", err)
		return
	}
	if Id != 2 {
		t.Errorf("Job.Insert() got = %d, want = %d", Id, 2)
		return
	}
	fmt.Printf("%+v\n", e)

	e = &Job{}
	gots, err := e.GetList()
	if err != nil {
		t.Errorf("Job.GetList() error = %v", err)
		return
	}
	if len(gots) != 2 {
		t.Errorf("Job.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}

	// Update
	e = NewJob("UNINSTALL", "test-1", "")
	want = e
	got, err = e.Update(2)
	if err != nil {
		t.Errorf("Job.Update() error = %v", err)
		return
	}
	if got.Uuid != want.Uuid || got.Status != want.Status ||
		got.Method != want.Method || got.Creator != want.Creator || got.ClusterId != want.ClusterId {
		t.Errorf("Job.Update() got = %+v, want = %+v", got, e)
		return
	}

	e = &Job{}
	e.ID = 2
	want = &got
	got, err = e.Get()
	if err != nil {
		t.Errorf("Job.Get() error = %v", err)
		return
	}
	if got.Uuid != want.Uuid || got.Status != want.Status ||
		got.Method != want.Method || got.Creator != want.Creator || got.ClusterId != want.ClusterId {
		t.Errorf("Job.Get() got = %v, wat = %v", got, want)
		return
	}

	e = &Job{}
	success, err := e.DeleteById(1)
	if err != nil {
		t.Errorf("Job.Delete() error = %v", err)
		return
	}
	if !success {
		t.Errorf("Job.Delete() success = %v, wat = %v", success, true)
		return
	}

	e = &Job{}
	gots, err = e.GetList()
	if err != nil {
		t.Errorf("Job.GetList() error = %v", err)
		return
	}
	if len(gots) != 1 {
		t.Errorf("Job.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}
}

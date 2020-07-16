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

	"github.com/jinzhu/gorm"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
)

func TestCluster(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Cluster{}
	e.DropTable()
	e.InitTable()
	// Drop Table
	//defer e.DropTable()

	//Insert
	e, err := NewCluster("fate-9999", "fate-9999", "fate", "v1.5.0", "")
	if err != nil {
		t.Errorf("Cluster.NewCluster() error = %v", err)
		return
	}
	e.ChartValues = map[string]interface{}{"Name": "fate-9999", "NameSpace": "fate-9999"}
	Id, err := e.Insert()
	if err != nil {
		t.Errorf("Cluster.Insert() error = %v", err)
		return
	}
	if Id < 1 {
		t.Errorf("Cluster.Insert() got = %d, want Id > 0", Id)
		return
	}
	// repeat insert
	Id, err = e.Insert()
	if err == nil {
		t.Error("Cluster.Insert() error = repeat insert")
		return
	}

	want := e
	e = &Cluster{Uuid: e.Uuid}
	got, err := e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.Status != want.Status ||
		got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Get() where Name=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Cluster{Name: "fate-9999"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Get() where Name=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Cluster{NameSpace: "fate-9999"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Get() where NameSpace=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Cluster{ChartName: "fate"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Get() where ChartName=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &Cluster{ChartVersion: "v1.5.0"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Get() where ChartVersion=fate-9999  got = %v, wat = %v", got, want)
		return
	}

	// Insert
	e, err = NewCluster("fate-10000", "fate-10000", "fate", "v1.4.0", "")
	if err != nil {
		t.Errorf("Cluster.NewCluster() error = %v", err)
		return
	}
	Id, err = e.Insert()
	if err != nil {
		t.Errorf("Cluster.Insert() error = %v", err)
		return
	}
	if Id != 2 {
		t.Errorf("Cluster.Insert() got = %d, want = %d", Id, 2)
		return
	}

	e = &Cluster{}
	gots, err := e.GetList()
	if err != nil {
		t.Errorf("Cluster.GetList() error = %v", err)
		return
	}
	if len(gots) != 2 {
		t.Errorf("Cluster.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}

	// Update
	e, err = NewCluster("fate-10001", "fate-10001", "fate-serving", "v1.4.1", "")
	if err != nil {
		t.Errorf("Cluster.NewCluster() error = %v", err)
		return
	}
	want = e
	got, err = e.Update(2)
	if err != nil {
		t.Errorf("Cluster.Update() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Update() got = %+v, want = %+v", got, e)
		return
	}

	e = &Cluster{}
	e.ID = 2
	want = &got
	got, err = e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.NameSpace != want.NameSpace ||
		got.Uuid != want.Uuid || got.ChartName != want.ChartName || got.ChartVersion != want.ChartVersion {
		t.Errorf("Cluster.Get() got = %v, wat = %v", got, want)
		return
	}

	e = &Cluster{Model: gorm.Model{ID: 1}}
	success, err := e.Delete()
	if err != nil {
		t.Errorf("Cluster.Delete() error = %v", err)
		return
	}
	if !success {
		t.Errorf("Cluster.Delete() success = %v, wat = %v", success, true)
		return
	}

	e = &Cluster{}
	gots, err = e.GetList()
	if err != nil {
		t.Errorf("Cluster.GetList() error = %v", err)
		return
	}
	if len(gots) != 1 {
		t.Errorf("Cluster.GetList() len(got) = %d, want = %d", len(gots), 1)
		return
	}

	e = &Cluster{}
	got, err = e.Get()
	if err != nil {
		t.Errorf("Cluster.Get() error = %v", err)
		return
	}
	err = got.SetStatus(ClusterStatusRunning)
	if err != nil {
		t.Errorf("Cluster.SetStatus() error = %v", err)
		return
	}
	fmt.Printf("%+v\n", got)
	got, err = e.Get()
	if got.Status != ClusterStatusRunning {
		t.Errorf("Cluster.Get() got.Status = %d, want = %s", got.Status, ClusterStatusRunning)
		return
	}
}

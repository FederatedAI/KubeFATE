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
	"helm.sh/helm/v3/pkg/chart"
)

func TestHelmChart(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &HelmChart{}

	e.InitTable()
	// Drop Table
	defer e.DropTable()

	//Insert
	templates := []*chart.File{
		{
			Name: "a.yaml",
			Data: []byte(`kind: ConfigMap
			apiVersion: v1
			metadata: `),
		},
		{
			Name: "b.yaml",
			Data: []byte(`kind: ConfigMap
			apiVersion: v1
			metadata: `),
		},
	}
	e = NewHelmChart("fate", "fate-chart", "fate-values", templates, "v1.5.0", "v1.5.0")
	Id, err := e.Insert()
	if err != nil {
		t.Errorf("HelmChart.Insert() error = %v", err)
		return
	}
	if Id != 1 {
		t.Errorf("HelmChart.Insert() got = %d, want = %d", Id, 1)
		return
	}
	// repeat insert
	Id, err = e.Insert()
	if err == nil {
		t.Error("HelmChart.Insert() error = repeat insert")
		return
	}

	want := e
	e = &HelmChart{Uuid: e.Uuid}
	got, err := e.Get()
	if err != nil {
		t.Errorf("HelmChart.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Get() where Name=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &HelmChart{Name: "fate"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("HelmChart.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Get() where Name=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &HelmChart{Chart: "fate-chart"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("HelmChart.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Get() where NameSpace=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &HelmChart{Values: "fate-values"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("HelmChart.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Get() where ChartName=fate-9999  got = %v, wat = %v", got, want)
		return
	}
	e = &HelmChart{Version: "v1.5.0"}
	got, err = e.Get()
	if err != nil {
		t.Errorf("HelmChart.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Get() where ChartVersion=fate-9999  got = %v, wat = %v", got, want)
		return
	}

	// Insert
	e = NewHelmChart("fate-1", "fate-chart-1", "fate-values-1", templates, "v1.5.1", "v1.5.1")
	Id, err = e.Insert()
	if err != nil {
		t.Errorf("HelmChart.Insert() error = %v", err)
		return
	}
	if Id != 2 {
		t.Errorf("HelmChart.Insert() got = %d, want = %d", Id, 2)
		return
	}

	e = &HelmChart{}
	gots, err := e.GetList()
	if err != nil {
		t.Errorf("HelmChart.GetList() error = %v", err)
		return
	}
	if len(gots) != 2 {
		t.Errorf("HelmChart.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}

	// Update
	e = NewHelmChart("fate-2", "fate-chart-2", "fate-values-2", templates, "v1.5.1-a", "v1.5.1")
	want = e
	got, err = e.Update(2)
	if err != nil {
		t.Errorf("HelmChart.Update() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Update() got = %+v, want = %+v", got, e)
		return
	}

	e = &HelmChart{}
	e.ID = 2
	want = &got
	got, err = e.Get()
	if err != nil {
		t.Errorf("HelmChart.Get() error = %v", err)
		return
	}
	if got.Name != want.Name || got.Chart != want.Chart ||
		got.Uuid != want.Uuid || got.Version != want.Version ||
		got.Values != want.Values {
		t.Errorf("HelmChart.Get() got = %v, wat = %v", got, want)
		return
	}

	e = &HelmChart{}
	success, err := e.Delete(1)
	if err != nil {
		t.Errorf("HelmChart.Delete() error = %v", err)
		return
	}
	if !success {
		t.Errorf("HelmChart.Delete() success = %v, wat = %v", success, true)
		return
	}

	e = &HelmChart{}
	gots, err = e.GetList()
	if err != nil {
		t.Errorf("HelmChart.GetList() error = %v", err)
		return
	}
	if len(gots) != 1 {
		t.Errorf("HelmChart.GetList() len(got) = %d, want = %d", len(gots), 2)
		return
	}
}

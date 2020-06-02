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
	"testing"
)

func TestNewFateCluster(t *testing.T) {
	InitConfigForTest()
	fate := NewCluster("fate-cluster1", "fate-nameSpaces", "fate", "v1.3.0")
	clusterUuid, error := Save(fate)
	if error == nil {
		t.Log("uuid: ", clusterUuid)
		clusterJustAddedUuid = clusterUuid
	}
}

func TestFindFateCluster(t *testing.T) {
	InitConfigForTest()
	fate := &Cluster{}
	results, error := Find(fate)
	if error == nil {
		t.Log(ToJson(results))
	}
}

func TestFindFateClusterByUuid(t *testing.T) {
	InitConfigForTest()
	t.Log("Find cluster just add: " + clusterJustAddedUuid)
	fate := &Cluster{}
	result, error := FindByUUID(fate, clusterJustAddedUuid)
	if error == nil {
		t.Log(ToJson(result))
		t.Log(result.(Cluster).Name)
	}
}

func TestDeleteClusterByUUID(t *testing.T) {
	InitConfigForTest()
	fate := &Cluster{}
	DeleteByUUID(fate, clusterJustAddedUuid)
}

//
//func TestReturnMethods(t *testing.T) {
//	InitConfigForTest()
//	fate := &Cluster{}
//	results, error := Find(fate)
//	if error == nil {
//		for _, v := range results {
//			oneFate := v.(Cluster)
//			t.Log(oneFate.GetUuid())
//			t.Log(oneFate.Name)
//		}
//	}
//}

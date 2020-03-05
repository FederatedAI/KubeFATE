package db

import (
	"testing"
)



func TestNewFateCluster(t *testing.T) {
	InitConfigForTest()
	fate := NewCluster("fate-cluster1", "fate-nameSpaces")
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

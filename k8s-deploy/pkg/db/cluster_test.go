package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

var clusterJustAddedUuid string

func TestNewCluster(t *testing.T) {
	InitConfigForTest()
	fate := NewCluster("fate-cluster1", "fate-nameSpaces","v1.3.0")
	clusterUuid, error := Save(fate)
	if error == nil {
		t.Log("uuid: ", clusterUuid)
		clusterJustAddedUuid = clusterUuid
	}
}

func TestFindCluster(t *testing.T) {
	InitConfigForTest()
	fate := &Cluster{}
	results, error := Find(fate)
	if error == nil {
		t.Log(ToJson(results))
	}
}

func TestFindClusterByUuid(t *testing.T) {
	InitConfigForTest()
	clusterJustAddedUuid = "f3a366f5-bf97-4be2-b49a-2137fe84a38b"
	t.Log("Find cluster just add: " + clusterJustAddedUuid)
	fate := &Cluster{}
	result, error := FindByUUID(fate, clusterJustAddedUuid)
	if error == nil {
		t.Log(ToJson(result))
		t.Log(result.(Cluster).Name)
	}
}

func TestUpdateCluster(t *testing.T) {
	InitConfigForTest()
	t.Log("Update: " + clusterJustAddedUuid)
	fate := &Cluster{}
	result, error := FindByUUID(fate, clusterJustAddedUuid)
	if error == nil {
		fate2Update := result.(Cluster)
		fate2Update.Name = "fate-cluster2"
		fate2Update.NameSpace = "fate-nameSpaces"

		UpdateByUUID(&fate2Update, clusterJustAddedUuid)
	}

	result, error = FindByUUID(fate, clusterJustAddedUuid)
	if error == nil {
		t.Log(ToJson(result))
	}
}

func TestDeleteByUUID(t *testing.T) {
	InitConfigForTest()
	fate := &Cluster{}
	DeleteByUUID(fate, clusterJustAddedUuid)
}

func TestReturnMethods(t *testing.T) {
	InitConfigForTest()
	fate := &Cluster{}
	results, error := Find(fate)
	if error == nil {
		for _, v := range results {
			oneFate := v.(Cluster)
			t.Log(oneFate.GetUuid())
			t.Log(oneFate.Name)
		}
	}
}

func TestFindClusterFindByUUID(t *testing.T) {
	InitConfigForTest()
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *Cluster
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				uuid: "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test",
			args: args{
				uuid: "aa3f3e57-79e3-497e-9ecb-ad778b119da1",
			},
			want: &Cluster{
				Uuid:      "2f41aabe-1610-4e4a-bc1c-9b24e9f8ec11",
				Name:      "fate-8888",
				NameSpace: "fate-8888",
				Version:   1,
				Metadata:  map[string]interface{}{},
				Status:    Creating_c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterFindByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterFindByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterFindByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterDeleteByUUID(t *testing.T) {
	InitConfigForTest()
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				uuid: "aa3f3e57-79e3-497e-9ecb-ad778b119da1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ClusterDeleteByUUID(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ClusterDeleteByUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindClusterList(t *testing.T) {
	InitConfigForTest()
	type args struct {
		args string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Cluster
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				args: "",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindClusterList(tt.args.args, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindClusterList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindClusterList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterDeleteAll(t *testing.T) {
	InitConfigForTest()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := ConnectDb()
	if err != nil {
		log.Error().Err(err).Msg("ConnectDb")
	}
	collection := db.Collection(new(Cluster).getCollection())
	filter := bson.D{}
	r, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("DeleteMany")
	}

	if r != nil && r.DeletedCount == 0 {
		log.Error().Msg("this record may not exist(DeletedCount==0)")
	}
	fmt.Println(r)
	return
}

func TestClusterFindByName(t *testing.T) {
	InitConfigForTest()
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    *Cluster
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterFindByName(tt.args.name, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterFindByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterFindByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_IsExisted(t *testing.T) {
	InitConfigForTest()
	type fields struct {
		Uuid             string
		Name             string
		NameSpace        string
		Version          int
		ChartVersion     string
		ChartValues      map[string]interface{}
		Values           string
		ChartName        string
		Type             string
		Metadata         map[string]interface{}
		Status           ClusterStatus
		Backend          ComputingBackend
		BootstrapParties Party
	}
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "",
			fields: fields{},
			args: args{
				name:      "fate-10000",
				namespace: "fate-10000",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				Version:      tt.fields.Version,
				ChartVersion: tt.fields.ChartVersion,
				ChartValues:  tt.fields.ChartValues,
				Values:       tt.fields.Values,
				ChartName:    tt.fields.ChartName,
				Type:         tt.fields.Type,
				Metadata:     tt.fields.Metadata,
				Status:       tt.fields.Status,
			}
			if got := cluster.IsExisted(tt.args.name, tt.args.namespace); got != tt.want {
				t.Errorf("Cluster.IsExisted() = %v, want %v", got, tt.want)
			}
		})
	}
}

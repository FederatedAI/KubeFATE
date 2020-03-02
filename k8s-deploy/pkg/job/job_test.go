// job
package job

import (
	"encoding/json"
	"fate-cloud-agent/pkg/db"
	"fate-cloud-agent/pkg/utils/config"
	"fate-cloud-agent/pkg/utils/logging"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
	logging.InitLog()
}

func TestMsa(t *testing.T) {

	d := ClusterArgs{
		Name:      "fate-10000",
		Namespace: "fate-10000",
		Version:   "V1.2.0",
		Data:      []byte(`{"egg":{"count":3},"exchange":{"ip":"192.168.1.1","port":9370},"modules":["proxy","egg","fateboard","fateflow","federation","metaService","mysql","redis","roll","python"],"partyId":10000,"proxy":{"nodePort":30010,"type":"NodePort"}}`),
	}
	b, err := json.Marshal(d)
	if err != nil {
		log.Err(err).Msg("err")
	}

	fmt.Printf("%s", b)

}

func TestClusterInstall(t *testing.T) {
	InitConfigForTest()
	type args struct {
		clusterArgs *ClusterArgs
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Job
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "job install fate-8888",
			args: args{
				clusterArgs: &ClusterArgs{
					Name:      "fate-8888",
					Namespace: "fate-8888",
					Version:   "v1.2.0",
					Data:      []byte(`{ "partyId":8888,"endpoint": { "ip":"10.184.111.187","port":30008}}`),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterInstall(tt.args.clusterArgs,"test")
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterInstall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterInstall() = %v, want %v", got, tt.want)
			}
			time.Sleep(60 * time.Second)
		})
	}
}

func TestClusterDelete(t *testing.T) {
	InitConfigForTest()
	type args struct {
		clusterId string
	}
	tests := []struct {
		name string
		args args
		want *db.Job
	}{
		// TODO: Add test cases.
		{
			name: "delete",
			args: args{
				clusterId: "5029628c-8886-4907-bced-6dbe3553c7ef",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got,_ := ClusterDelete(tt.args.clusterId,"test"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterDelete() = %v, want %v", got, tt.want)
			}
			time.Sleep(30 * time.Second)
		})
	}
}

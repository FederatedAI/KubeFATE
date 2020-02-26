package cli

import (
	"fate-cloud-agent/pkg/utils/logging"
	"os"
	"reflect"
	"testing"
)

func TestCluster(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		Args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"cluster -help",
			args{[]string{os.Args[0], "cluster", "--help"}},
		},
		{
			"cluster list",
			args{[]string{os.Args[0], "cluster", "list"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.Args)
		})
	}
}

func Test_yamlToJson(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{
				bytes: []byte(`partyId: 10000
endpoint:
  ip: 10.184.111.187
  port: 30000`),
			},
			want:    []byte(`{"endpoint":{"ip":"10.184.111.187","port":30000},"partyId":10000}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := yamlToJson(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("yamlToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("yamlToJson() = %s, want %s", got, tt.want)
			}
		})
	}
}

package service

import (
	"fate-cloud-agent/pkg/utils/config"
	"fate-cloud-agent/pkg/utils/logging"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestInstall(t *testing.T) {

	_ = config.InitViper()
	viper.AddConfigPath("../../")
	_ = viper.ReadInConfig()
	logging.InitLog()
	err := RepoAddAndUpdate()
	if err != nil {
		panic(err)
	}
	type args struct {
		namespace string
		name      string
		version   string
		value     *Value
	}
	tests := []struct {
		name    string
		args    args
		want    *releaseElement
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "install fate",
			args: args{
				namespace: "fate-10000",
				name:      "fate-10000",
				version:   "v1.2.0",
				value:     &Value{Val: []byte(`{ "partyId":10000,"endpoint": { "ip":"10.184.111.187","port":30000}}`), T: "json"},
			},
			want: &releaseElement{
				Name:       "fate",
				Namespace:  "fate",
				Revision:   "1",
				Status:     "deployed",
				Chart:      "fate-1.2.0",
				AppVersion: "1.2.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Install(tt.args.namespace, tt.args.name, tt.args.version, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Install() = %v, want %v", got, tt.want)
			}

		})
	}
}

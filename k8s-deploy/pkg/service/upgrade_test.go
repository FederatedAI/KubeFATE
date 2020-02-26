package service

import (
	"fate-cloud-agent/pkg/utils/config"
	"fate-cloud-agent/pkg/utils/logging"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestUpgrade(t *testing.T) {
	_ = config.InitViper()
	viper.AddConfigPath("../../")
	_ = viper.ReadInConfig()
	logging.InitLog()
	_ = os.Setenv("FATECLOUD_CHART_PATH", "../../")
	type args struct {
		namespace string
		name      string
		version   string
		value     *Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{
				namespace: "fate-10000",
				name:      "fate-10000",
				version:   "v1.2.0",
				value:     &Value{
					Val: []byte(`{ "partyId":10000,"endpoint": { "ip":"10.184.111.187","port":30030}}`),
					T:   "json",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Upgrade(tt.args.namespace, tt.args.name, tt.args.version, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Upgrade() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

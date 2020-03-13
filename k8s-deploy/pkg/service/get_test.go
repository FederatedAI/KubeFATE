package service

import (
	"fate-cloud-agent/pkg/utils/config"
	"fate-cloud-agent/pkg/utils/logging"
	"github.com/spf13/viper"
	"testing"
)

func TestGet(t *testing.T) {
	_ = config.InitViper()
	viper.AddConfigPath("../../")
	_ = viper.ReadInConfig()
	logging.InitLog()
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.

		{
			name: "fate name no find",
			args: args{
				namespace: "fate-10001",
				name:      "fate-10001",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fate namespace no find",
			args: args{
				namespace: "fate-10001",
				name:      "fate-10000",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fate",
			args: args{
				namespace: "fate-10000",
				name:      "fate-10000",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"github.com/spf13/viper"
	"reflect"
	"testing"

	"helm.sh/helm/v3/pkg/release"
)

func TestDelete(t *testing.T) {

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
		want    *release.UninstallReleaseResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete",
			args: args{
				namespace: "fate-10000",
				name:      "fate-10000",
			},
			want: &release.UninstallReleaseResponse{
				Release: nil,
				Info:    "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Delete(tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Info, tt.want.Info) {
				t.Errorf("Delete() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

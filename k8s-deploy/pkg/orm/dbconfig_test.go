package orm

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"testing"
)

func TestGetDbConfig(t *testing.T) {
	InitConfigForTest()
	tests := []struct {
		name string
		want *DbConfig
	}{
		// TODO: Add test cases.
		{
			name: "Environment variables as configuration files",
			want: &DbConfig{
				DbType:   os.Getenv("FATECLOUD_DB_TYPE"),
				Host:     os.Getenv("FATECLOUD_DB_HOST"),
				Port:     os.Getenv("FATECLOUD_DB_PORT"),
				Name:     os.Getenv("FATECLOUD_DB_NAME"),
				Username: os.Getenv("FATECLOUD_DB_USERNAME"),
				Password: os.Getenv("FATECLOUD_DB_PASSWORD"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDbConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDbConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
}

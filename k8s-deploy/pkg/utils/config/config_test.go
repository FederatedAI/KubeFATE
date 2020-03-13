package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestConfig_DirExists(t *testing.T) {
	tmpDir := os.TempDir()
	exists := DirExists(tmpDir)
	if exists != true {
		t.Errorf("%s exists but DirExists return false \n", tmpDir)
	}

	// construct a random dir path
	tmpDir = time.Now().Format(time.RFC3339Nano)
	exists = DirExists(tmpDir)
	if exists != false {
		t.Errorf("%s does not exist but DirExists return true \n", tmpDir)
	}
}

func TestConfig_InitViper(t *testing.T) {
	_ = InitViper()

	viper.AddConfigPath("../../../")

	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("Fatal error config file: %s \n", err)
	}

	result := viper.Get("mongo")
	if result == "" {
		t.Errorf("Can not read mongo")
	}

	t.Log(result)
	result = viper.Get("mongo.url")
	t.Log(result)
}

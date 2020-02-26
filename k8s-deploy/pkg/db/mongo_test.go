package db

import (
	"context"
	"fate-cloud-agent/pkg/utils/config"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestDB_InitKubeFATEDatabase(t *testing.T) {
	InitConfigForTest()

	initKubeFATEDatabase()

	// Log the constructed mongo url
	t.Log(kubeFateDB.returnMongoInfo())

	// Log the constructed mongo url after env was changed
	os.Setenv("FATECLOUD_MONGO_USERNAME", "test")
	os.Setenv("FATECLOUD_MONGO_PASSWORD", "test")

	// Sleep for a while
	time.Sleep(2 * time.Second)

	// Log the constructed mongo url
	t.Log(kubeFateDB.returnMongoInfo())

	expectedResult := "mongodb://test:test@localhost:27017"
	if kubeFateDB.returnMongoInfo() != expectedResult {
		t.Errorf("Expecte %s but get %s", expectedResult, kubeFateDB.returnMongoInfo())
	}

	// Log the constructed mongo url after env was changed
	os.Setenv("FATECLOUD_MONGO_USERNAME", "")
	os.Setenv("FATECLOUD_MONGO_PASSWORD", "")
}

func TestDB_ConnectDb(t *testing.T) {
	InitConfigForTest()

	db, _ := ConnectDb()
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)

	err := db.Client().Ping(ctx, readpref.Primary())

	if err != nil {
		t.Errorf("Unable to ping db: %s", err)
	}
}

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
}

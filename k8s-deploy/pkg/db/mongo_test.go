/*
* Copyright 2019-2020 VMware, Inc.
* 
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
* 
*/
package db

import (
	"context"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

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

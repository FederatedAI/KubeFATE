/*
 * Copyright 2019-2021 VMware, Inc.
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

package api

import (
	"fmt"
	"os"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

func initUser() error {
	err := generateAdminUser()
	if err != nil {
		return fmt.Errorf("generate admin user error: %s ", err.Error())
	}
	return nil
}

func initDb() error {
	var err error

	for i := 0; i < 3; i++ {
		err = orm.InitDB()
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("initialization failed: %s", err)
}

func initTables() error {
	err := new(modules.User).InitTable()
	if err != nil {
		return err
	}
	err = new(modules.Cluster).InitTable()
	if err != nil {
		return err
	}
	err = new(modules.HelmChart).InitTable()
	if err != nil {
		return err
	}
	err = new(modules.Job).InitTable()
	if err != nil {
		return err
	}
	return nil
}

// Run starts the API server
// @title KubeFATE service API
// @version v1
// @description This is a KubeFATE.
// @contact.name API Support
// @contact.url https://github.com/FederatedAI/KubeFATE
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Run() error {
	log.Info().Msgf("logLevel: %v", viper.GetString("log.level"))
	log.Info().Msgf("api version: %v", APIVersion)
	log.Info().Msgf("service version: %v", ServiceVersion)
	log.Info().Msgf("DbType: %v", viper.GetString("db.type"))
	log.Info().Msgf("LogNocolor: %v", viper.GetString("log.nocolor"))
	log.Info().Msgf("server: [%s:%s]", viper.GetString("server.address"), viper.GetString("server.port"))

	err := initDb()
	if err != nil {
		log.Error().Err(err).Msg("initDb error, ")
		return err
	}

	modules.DB = orm.DB

	log.Info().Msg("Database connection Successful")

	err = initTables()
	if err != nil {
		log.Error().Err(err).Msg("initTables error, ")
		return err
	}

	log.Info().Msg("Table initialization succeeded")

	err = initUser()
	if err != nil {
		log.Error().Err(err).Msg("initUser error, ")
		return err
	}

	log.Info().Msg("User created Successfully")

	// use gin.New() instead
	r := gin.New()

	// use default recovery
	r.Use(gin.Recovery())

	// reset caller info level to identify http server log from normal log
	// customizedLog := log.With().CallerWithSkipFrameCount(9).Logger()
	// use customized logger
	r.Use(
		logger.SetLogger(
			logger.WithUTC(true),
			logger.WithLogger(logging.GetGinLogger),
		),
	)

	Router(r)

	address := viper.GetString("server.address")
	port := viper.GetString("server.port")
	endpoint := fmt.Sprintf("%s:%s", address, port)

	// It is weird that release mode won't output serving info
	if os.Getenv("GIN_MODE") == "release" {
		log.Info().Msg("Listening and serving HTTP on " + address + ":" + port)
	}

	log.Info().Msg("Gin configuration Successful")

	err = r.Run(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("gin run error, ")
		return err
	}
	return nil
}

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

package orm

import (
	"fmt"
	"gorm.io/gorm/logger"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Database interface {
	Open(logLevel logger.LogLevel) (db *gorm.DB, err error)
}

func getDbType(Type string) (Database, error) {
	var database Database
	switch Type {
	case "mysql":
		database = new(Mysql)
	case "sqlite":
		database = new(Sqlite)
	default:
		err := fmt.Errorf("unknown database type: %s, please use 'mysql' or 'sqlite'", Type)
		log.Error().Str("Type", Type).Err(err).Msg("unknown db type")
		return nil, err
	}
	return database, nil
}

func getLogLevel(Type string) logger.LogLevel {
	switch Type {
	case "debug":
		return logger.Info
	default:
		return logger.Silent
	}
}

func Setup() (*gorm.DB, error) {
	database, err := getDbType(viper.GetString("db.type"))
	if err != nil {
		return nil, err
	}
	logLevel := getLogLevel(viper.GetString("log.level"))
	return database.Open(logLevel)
}

var DB *gorm.DB

func InitDB() error {
	db, err := Setup()
	if err != nil {
		return err
	}
	DB = db
	return nil
}

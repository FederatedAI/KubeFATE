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
	"bytes"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

func getDbConfig() *DbConfig {
	return &DbConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Name:     viper.GetString("db.name"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
	}
}

func (e *Mysql) Open(logLevel logger.LogLevel) (db *gorm.DB, err error) {
	dbConfig := getDbConfig()
	var conn bytes.Buffer
	conn.WriteString(dbConfig.Username)
	conn.WriteString(":")
	conn.WriteString(dbConfig.Password)
	conn.WriteString("@tcp(")
	conn.WriteString(dbConfig.Host)
	conn.WriteString(":")
	conn.WriteString(dbConfig.Port)
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(dbConfig.Name)
	conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=10s")
	dsn := conn.String()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      false,
		},
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
}

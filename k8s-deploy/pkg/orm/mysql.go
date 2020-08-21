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

package orm

import (
	"bytes"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/zerolog/log"
)

var DBCLIENT *gorm.DB

type Mysql struct {
}

func (e *Mysql) GetConnect() string {
	dbConfig := GetDbConfig()
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
	return conn.String()
}

func (e *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

func (e *Mysql) Setup() error {

	var err error
	var db Database

	db = new(Mysql)
	mysqlConn := db.GetConnect()
	log.Info().Msg(mysqlConn)

	DBCLIENT, err = db.Open("mysql", mysqlConn)

	if err != nil {
		log.Error().Msgf("%s connect error %v", "mysql", err)
		return err
	} else {
		log.Info().Msgf("%s connect success!", "mysql")
	}
	return nil
}

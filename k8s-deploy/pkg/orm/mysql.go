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
	conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=1000ms")
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
		log.Fatal().Msgf("%s connect error %v", "mysql", err)
		return err
	} else {
		log.Info().Msgf("%s connect success!", "mysql")
	}
	return nil
}

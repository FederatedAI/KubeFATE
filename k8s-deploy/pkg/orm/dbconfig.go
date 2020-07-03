package orm

import "github.com/spf13/viper"

type DbConfig struct {
	DbType   string
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

func GetDbConfig() *DbConfig {
	return &DbConfig{
		DbType:   viper.GetString("db.type"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Name:     viper.GetString("db.name"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
	}
}


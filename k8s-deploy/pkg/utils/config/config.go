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
package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const cUserSpecifiedPATH string = "FATECLOUD_CONFIG_PATH"
const cEnvironmentPrefix string = "FATECLOUD"

// InitViper initial a viper instance
func InitViper() {
	setDefaultConfig()
	// For environment variable
	viper.SetEnvPrefix("FATECLOUD")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	path, _ := filepath.Abs(".")
	viper.AddConfigPath(path)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return
}

func setDefaultConfig() {
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", "8080")

	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.file", "file::memory:?cache=shared")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.nocolor", "false")

	viper.SetDefault("repo.name", "kubefate")
	viper.SetDefault("repo.url", "https://federatedai.github.io/KubeFATE")

	viper.SetDefault("user.username", "admin")
	viper.SetDefault("user.password", "admin")

	viper.SetDefault("serviceurl", "localhost:8080")
}

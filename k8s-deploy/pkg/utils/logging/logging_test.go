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
package logging

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"
)

func TestInitLog(t *testing.T) {
	logLevel := "info"
	viper.Set("log.nocolor", "true")
	viper.Set("log.level", logLevel)
	InitLog()
	if log.Logger.GetLevel().String() != logLevel {
		t.Errorf("log level Configuration error")
	}
	log.Trace().Msg("Trace")
	log.Info().Msg("Info")
	log.Debug().Msg("Debug")
	log.Warn().Msg("Warn")

}

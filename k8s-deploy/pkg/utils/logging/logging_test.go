package logging

import (
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

func TestInitLog(t *testing.T) {
	os.Setenv("FATECLOUD_LOG_LEVEL", "debug")
	InitLog()
	log.Info().Str("Logger Level", log.Logger.GetLevel().String())

	log.Debug().Msg("debug")
	log.Info().Msg("Info")

}

package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// InitLog initials a log instance with specified config
func InitLog() {
	logLevel := viper.GetString("log.level")


	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		// Info is the default log level
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			// TODO(JHC):
			// 1. output device should read from config file
			Out:        os.Stderr,
			NoColor:    false,
			TimeFormat: time.RFC3339,
		},
	)

	log.Logger = log.With().Caller().CallerWithSkipFrameCount(2).Logger()

}

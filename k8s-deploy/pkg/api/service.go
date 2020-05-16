package api

import (
	"fmt"
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func initUser() error {
	err := generateAdminUser()
	if err != nil {
		return fmt.Errorf("generate admin user error: %s\n", err)
	}
	return nil
}

// Run starts the API server
func Run() {
	log.Info().Msgf("logLevel: %v", viper.GetString("log.level"))
	log.Info().Msgf("api version: %v", ApiVersion)
	log.Info().Msgf("service version: %v", ServiceVersion)
	err := initUser()
	if err != nil {
		log.Error().Err(err).Msg("initUser error, ")
		return
	}
	//err := service.InitKubeConfig()
	//if err != nil {
	//	panic(err)
	//}

	// use gin.New() instead
	r := gin.New()

	// use default recovery
	r.Use(gin.Recovery())

	// reset caller info level to identify http server log from normal log
	customizedLog := log.With().CallerWithSkipFrameCount(9).Logger()
	// use customized logger
	r.Use(logger.SetLogger(logger.Config{
		Logger: &customizedLog,
		UTC:    true,
	}))

	Router(r)

	address := viper.GetString("server.address")
	port := viper.GetString("server.port")
	endpoint := fmt.Sprintf("%s:%s", address, port)

	// It is weird that release mode won't output serving info
	if os.Getenv("GIN_MODE") == "release" {
		log.Info().Msg("Listening and serving HTTP on " + address + ":" + port)
	}

	err = r.Run(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("gin run error, ")
		return
	}
}

package api

import (
	"fate-cloud-agent/pkg/service"
	"fmt"
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func initUser() {
	err := generateAdminUser()
	if err != nil {
		panic(fmt.Errorf("generate admin user error: %s\n", err))
	}
}

// Run starts the API server
func Run() {
	initUser()
	//err := service.InitKubeConfig()
	//if err != nil {
	//	panic(err)
	//}
	err := service.RepoAddAndUpdate()
	if err != nil {
		panic(err)
	}

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

	_ = r.Run(endpoint)
}

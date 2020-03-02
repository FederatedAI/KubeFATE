package main

import (
	"fate-cloud-agent/pkg/cli"
	"fate-cloud-agent/pkg/utils/config"
	"fate-cloud-agent/pkg/utils/logging"
	"fmt"
	"os"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(fmt.Errorf("Unable to read in configuration: %s\n", err))
	}

	logging.InitLog()

	cli.Run(os.Args)
}

package main

import (
	"fmt"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/cli"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"os"
)

func main() {
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Unable to read in configuration: %s\n", err)
		return
	}

	logging.InitLog()

	cli.Run(os.Args)
}

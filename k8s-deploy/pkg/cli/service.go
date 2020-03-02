package cli

import (
	"fate-cloud-agent/pkg/api"
	"github.com/urfave/cli/v2"
)

func serviceCommand() *cli.Command {
	return &cli.Command{
		Name:   "service",
		Usage:  "service",
		Action: serviceRun,
	}
}

func serviceRun(c *cli.Context) error {
	api.Run()
	return nil
}

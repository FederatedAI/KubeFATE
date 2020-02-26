package cli

import (
	"fate-cloud-agent/pkg/service"
	"github.com/urfave/cli/v2"
)

func configCommand() *cli.Command {
	return &cli.Command{
		Name:   "config",
		Usage:  "config",
		Action: conf,
	}
}

func conf(c *cli.Context) error {

	return service.Conf()
}

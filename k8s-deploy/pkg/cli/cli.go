package cli

import (
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"sort"
)

func initCommandLine() *cli.App {
	app := &cli.App{

		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			serviceCommand(),
			ClusterCommand(),
			JobCommand(),
			UserCommand(),
			ChartCommand(),
			VersionCommand(),
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	return app
}

func Run(Args []string) {
	app := initCommandLine()
	err := app.Run(Args)
	if err != nil {
		fmt.Println(err)
	}
}

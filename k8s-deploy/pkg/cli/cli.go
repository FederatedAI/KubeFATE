package cli

import (
	"github.com/urfave/cli/v2" // imports as package "cli"
	"log"
	"sort"
)

func initCommandLine() *cli.App {
	app := &cli.App{

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "n",
				Usage: "fate name",
			},
		},
		Commands: []*cli.Command{
			serviceCommand(),
			ClusterCommand(),
			JobCommand(),
			UserCommand(),
			configCommand(),
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
		log.Fatal(err)
	}
}

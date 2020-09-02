package cli

import (
	"github.com/urfave/cli/v2"
)

func NamespaceCommand() *cli.Command {
	return &cli.Command{
		Name:  "namespace",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			NamespaceListCommand(),
		},
		Usage: "List namespace",
	}
}

func NamespaceListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags:   []cli.Flag{},
		Usage:   "Show Namespace list",
		Action: func(c *cli.Context) error {
			Namespace := new(Namespace)
			return GetItemList(Namespace)
		},
	}
}

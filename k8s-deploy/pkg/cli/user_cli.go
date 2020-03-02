package cli

import (
	"errors"
	"github.com/urfave/cli/v2"
)

func UserCommand() *cli.Command {
	return &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
		},
		Subcommands: []*cli.Command{
			UserListCommand(),
			UserInfoCommand(),
		},
		Usage: "add a task to the list",
	}
}

func UserListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
		},
		Usage: "show job list",
		Action: func(c *cli.Context) error {
			User := new(User)
			return getItemList(User)
		},
	}
}

func UserInfoCommand() *cli.Command {
	return &cli.Command{
		Name: "describe",
		Flags: []cli.Flag{
		},
		Usage: "show User info",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			User := new(User)
			return getItem(User, uuid)
		},
	}
}

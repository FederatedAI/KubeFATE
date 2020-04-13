package cli

import (
	"errors"
	"github.com/urfave/cli/v2"
)

func UserCommand() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			UserListCommand(),
			UserInfoCommand(),
		},
		Usage: "List all users and describe a user's info",
	}
}

func UserListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags:   []cli.Flag{},
		Usage:   "List all users",
		Action: func(c *cli.Context) error {
			User := new(User)
			return GetItemList(User)
		},
	}
}

func UserInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "describe",
		Flags: []cli.Flag{},
		Usage: "Describe a user's detail info",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			User := new(User)
			return GetItem(User, uuid)
		},
	}
}

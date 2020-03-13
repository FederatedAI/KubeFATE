package cli

import (
	"errors"
	"github.com/urfave/cli/v2"
)

func JobCommand() *cli.Command {
	return &cli.Command{
		Name: "job",
		Flags: []cli.Flag{
		},
		Subcommands: []*cli.Command{
			JobListCommand(),
			JobInfoCommand(),
			JobDeleteCommand(),
		},
		Usage: "List jobs, describe and delete a job",
	}
}

func JobListCommand() *cli.Command {
	return &cli.Command{
		Name: "list",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
		},
		Usage: "Show job list",
		Action: func(c *cli.Context) error {
			cluster := new(Job)
			return getItemList(cluster)
		},
	}
}

func JobDeleteCommand() *cli.Command {
	return &cli.Command{
		Name: "delete",
		Aliases: []string{"del"},
		Flags: []cli.Flag{
		},
		Usage: "Delete a job",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			cluster := new(Job)
			return deleteItem(cluster, uuid)
		},
	}
}

func JobInfoCommand() *cli.Command {
	return &cli.Command{
		Name: "describe",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "uuid",
				Value: "",
				Usage: "Describe a job with given UUID",
			},
		},
		Usage: "Show job's details info",
		Action: func(c *cli.Context) error {

			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			Job := new(Job)
			return getItem(Job, uuid)
		},
	}
}

/*
* Copyright 2019-2020 VMware, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
 */
package cli

import (
	"errors"
	"github.com/urfave/cli/v2"
)

func JobCommand() *cli.Command {
	return &cli.Command{
		Name:  "job",
		Flags: []cli.Flag{},
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
		Name:    "list",
		Aliases: []string{"ls"},
		Flags:   []cli.Flag{},
		Usage:   "Show job list",
		Action: func(c *cli.Context) error {
			cluster := new(Job)
			return GetItemList(cluster)
		},
	}
}

func JobDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Flags:   []cli.Flag{},
		Usage:   "Delete a job",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			cluster := new(Job)
			return DeleteItem(cluster, uuid)
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
			return GetItem(Job, uuid)
		},
	}
}

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

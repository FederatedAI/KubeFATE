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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/job"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"sigs.k8s.io/yaml"

	"io/ioutil"
)

func ClusterCommand() *cli.Command {
	return &cli.Command{
		Name:  "cluster",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			ClusterListCommand(),
			ClusterInfoCommand(),
			ClusterDeleteCommand(),
			ClusterInstallCommand(),
			ClusterUpdateCommand(),
		},
		Usage: "Manage cluster install, delete and update",
	}
}

func ClusterListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"A"},
				Value:   false,
				Usage:   "List all clusters including deleted ones",
			},
		},
		Usage: "Show all clusters list",
		Action: func(c *cli.Context) error {
			all := c.Bool("all")
			cluster := new(Cluster)
			cluster.all = all
			log.Debug().Bool("all", all).Msg("all")
			return GetItemList(cluster)
		},
	}
}

func ClusterInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "describe",
		Flags: []cli.Flag{},
		Usage: "Describe a cluster's detail info",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			cluster := new(Cluster)
			return GetItem(cluster, uuid)
		},
	}
}

func ClusterDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Flags:   []cli.Flag{},
		Usage:   "Delete a cluster",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}

			cluster := new(Cluster)
			log.Debug().Str("uuid", uuid).Msg("cluster delete uuid")
			return DeleteItem(cluster, uuid)
		},
	}
}

func ClusterInstallCommand() *cli.Command {
	return &cli.Command{
		Name: "install",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Value:    "",
				Usage:    "chart cluster.yaml",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "cover",
				Value: false,
				Usage: "If the cluster already exists, overwrite the installation",
			},
		},
		Usage: "Install a new cluster",
		Action: func(c *cli.Context) error {

			valTemValPath := c.String("file")
			log.Debug().Str("file", valTemValPath).Msg("install flag")
			cover := c.Bool("cover")
			log.Debug().Bool("cover", cover).Msg("install flag")

			clusterConfig, err := ioutil.ReadFile(valTemValPath)
			if err != nil {
				return err
			}
			log.Debug().Str("yaml", string(clusterConfig)).Msg("ReadFile success")

			var m map[string]interface{}
			err = yaml.Unmarshal(clusterConfig, &m)
			if err != nil {
				return err
			}

			name, ok := m["name"]
			if !ok {
				return errors.New("name not found, check your cluster file")
			}

			namespace, ok := m["namespace"]
			if !ok {
				return errors.New("namespace not found, check your cluster file")
			}

			chartVersion, ok := m["chartVersion"]
			if !ok {
				return errors.New("chartVersion not found, check your cluster file")
			}

			chartName, ok := m["chartName"]
			if !ok {
				chartName = ""
				//return errors.New("chartName not found, check your cluster file")
			}

			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			valBJ, err := json.Marshal(m)

			if err != nil {
				return err
			}

			cluster := new(Cluster)
			args := &job.ClusterArgs{
				Name:         name.(string),
				Namespace:    namespace.(string),
				ChartName:    chartName.(string),
				ChartVersion: chartVersion.(string),
				Cover:        cover,
				Data:         valBJ,
			}

			body, err := json.Marshal(args)
			if err != nil {
				return err
			}
			return PostItem(cluster, body)
		},
	}
}

func ClusterUpdateCommand() *cli.Command {
	return &cli.Command{
		Name: "update",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Value:    "",
				Usage:    "chart cluster.yaml",
				Required: true,
			},
		},
		Usage: "Update a cluster",
		Action: func(c *cli.Context) error {
			valTemValPath := c.String("file")

			clusterConfig, err := ioutil.ReadFile(valTemValPath)
			if err != nil {
				return err
			}

			log.Debug().Str("yaml", string(clusterConfig)).Msg("ReadFile success")

			var m map[string]interface{}
			err = yaml.Unmarshal(clusterConfig, &m)
			if err != nil {
				return err
			}

			name, ok := m["name"]
			if !ok {
				return errors.New("name not found, check your cluster file")
			}

			namespace, ok := m["namespace"]
			if !ok {
				return errors.New("namespace not found, check your cluster file")
			}

			chartVersion, ok := m["chartVersion"]
			if !ok {
				return errors.New("chartVersion not found, check your cluster file")
			}

			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			valBJ, err := json.Marshal(m)

			if err != nil {
				return err
			}

			cluster := new(Cluster)
			args := &job.ClusterArgs{
				Name:         name.(string),
				Namespace:    namespace.(string),
				ChartVersion: chartVersion.(string),
				Cover:        false,
				Data:         valBJ,
			}

			body, err := json.Marshal(args)
			if err != nil {
				return err
			}
			return PutItem(cluster, body)
		},
	}
}

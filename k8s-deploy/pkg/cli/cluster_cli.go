package cli

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"sigs.k8s.io/yaml"

	"io/ioutil"
)

func ClusterCommand() *cli.Command {
	return &cli.Command{
		Name: "cluster",
		Flags: []cli.Flag{
		},
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
				Usage:   "chart valTemVal.yaml",
			},
		},
		Usage: "show cluster list",
		Action: func(c *cli.Context) error {
			all := c.Bool("all")
			cluster := new(Cluster)
			cluster.all = all
			log.Debug().Bool("all",all).Msg("all")
			return getItemList(cluster)
		},
	}
}

func ClusterInfoCommand() *cli.Command {
	return &cli.Command{
		Name: "describe",
		Flags: []cli.Flag{
		},
		Usage: "show cluster info",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}
			cluster := new(Cluster)
			return getItem(cluster, uuid)
		},
	}
}

func ClusterDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Flags: []cli.Flag{
		},
		Usage: "cluster delete",
		Action: func(c *cli.Context) error {
			var uuid string
			if c.Args().Len() > 0 {
				uuid = c.Args().Get(0)
			} else {
				return errors.New("not uuid")
			}

			cluster := new(Cluster)
			log.Debug().Str("uuid", uuid).Msg("cluster delete uuid")
			return deleteItem(cluster, uuid)
		},
	}
}

func ClusterInstallCommand() *cli.Command {
	return &cli.Command{
		Name: "install",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "chart valTemVal.yaml",
			},
		},
		Usage: "cluster delete",
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
			delete(m, "name")

			namespace, ok := m["namespace"]
			if !ok {
				return errors.New("namespace not found, check your cluster file")
			}
			delete(m, "namespace")

			version, ok := m["version"]
			if !ok {
				return errors.New("version not found, check your cluster file")
			}
			delete(m, "version")

			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			valBJ, err := json.Marshal(m)

			if err != nil {
				return err
			}

			cluster := new(Cluster)
			args := struct {
				Name      string
				Namespace string
				Version   string
				Data      []byte
			}{
				Namespace: namespace.(string),
				Name:      name.(string),
				Version:   version.(string),
				Data:      valBJ,
			}

			body, err := json.Marshal(args)
			if err != nil {
				return err
			}
			return postItem(cluster, body)
		},
	}
}

func ClusterUpdateCommand() *cli.Command {
	return &cli.Command{
		Name: "update",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "chart valTemVal.yaml",
			},
		},
		Usage: "cluster Upgrade",
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
			delete(m, "name")

			namespace, ok := m["namespace"]
			if !ok {
				return errors.New("namespace not found, check your cluster file")
			}
			delete(m, "namespace")

			version, ok := m["version"]
			if !ok {
				return errors.New("version not found, check your cluster file")
			}
			delete(m, "version")

			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			valBJ, err := json.Marshal(m)

			if err != nil {
				return err
			}

			cluster := new(Cluster)
			args := struct {
				Name      string
				Namespace string
				Version   string
				Data      []byte
			}{
				Namespace: namespace.(string),
				Name:      name.(string),
				Version:   version.(string),
				Data:      valBJ,
			}

			body, err := json.Marshal(args)
			if err != nil {
				return err
			}
			return putItem(cluster, body)
		},
	}
}

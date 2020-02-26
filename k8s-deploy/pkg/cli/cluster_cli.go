package cli

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"sigs.k8s.io/yaml"

	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/rand"
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
			ClusterUpgradeCommand(),
		},
		Usage: "add a task to the list",
	}
}

func ClusterListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
		},
		Usage: "show cluster list",
		Action: func(c *cli.Context) error {
			cluster := new(Cluster)
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
			&cli.StringFlag{
				Name:    "namespace",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "k8s namespace",
			},
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Value:   "",
				Usage:   "chart version",
			},
		},
		Usage: "cluster delete",
		Action: func(c *cli.Context) error {

			valTemValPath := c.String("file")

			valBY, err := ioutil.ReadFile(valTemValPath)
			if err != nil {
				return err
			}
			log.Debug().Str("yaml",string(valBY)).Msg("ReadFile success")
			valBJ, err := yamlToJson(valBY)
			if err != nil {
				return err
			}
			var name string
			if c.Args().Len() > 0 {
				name = c.Args().Get(0)
			} else {
				name = "fate-" + rand.String(4)
			}

			cluster := new(Cluster)
			args := struct {
				Name      string
				Namespace string
				Version   string
				Data      []byte
			}{
				Namespace: c.String("namespace"),
				Name:      name,
				Version:   c.String("version"),
				Data:      valBJ,
			}
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			body, err := json.Marshal(args)
			if err != nil {
				return err
			}
			return postItem(cluster, body)
		},
	}
}

func ClusterUpgradeCommand() *cli.Command {
	return &cli.Command{
		Name: "upgrade",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "chart valTemVal.yaml",
			},
			&cli.StringFlag{
				Name:    "namespace",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "k8s namespace",
			},
		},
		Usage: "cluster Upgrade",
		Action: func(c *cli.Context) error {
			valTemValPath := c.String("file")

			valBY, err := ioutil.ReadFile(valTemValPath)
			if err != nil {
				return err
			}

			log.Debug().Str("yaml",string(valBY)).Msg("ReadFile success")

			valBJ, err := yamlToJson(valBY)
			if err != nil {
				return err
			}
			var name string
			if c.Args().Len() > 0 {
				name = c.Args().Get(0)
			} else {
				name = "fate-" + rand.String(4)
			}

			cluster := new(Cluster)
			args := struct {
				Name      string
				Namespace string
				Version   string
				Data      []byte
			}{
				Namespace: c.String("namespace"),
				Name:      name,
				Version:   "",
				Data:      valBJ,
			}
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			body, err := json.Marshal(args)
			if err != nil {
				return err
			}
			return putItem(cluster, body)
		},
	}
}

func yamlToJson(bytes []byte) ([]byte, error) {
	var m map[string]interface{}
	err := yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(m)

}

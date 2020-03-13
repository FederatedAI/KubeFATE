package service

import (
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"os"
)

func Delete(namespace, name string) (*release.UninstallReleaseResponse, error) {

	EnvCs.Lock()
	err := os.Setenv("HELM_NAMESPACE", namespace)
	if err != nil {
		panic(err)
	}
	settings := cli.New()
	EnvCs.Unlock()

	cfg := new(action.Configuration)
	client := action.NewUninstall(cfg)

	if err := cfg.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		return nil, err
	}

	res, err := client.Run(name)
	if err != nil {
		return nil, err
	}

	log.Debug().Interface("resInfo", res.Info).Msg("delete result")

	return res, nil
}

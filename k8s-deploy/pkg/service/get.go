package service

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"os"
)

func Get(namespace, name string) (*release.Release, error) {
	EnvCs.Lock()
	err := os.Setenv("HELM_NAMESPACE", namespace)
	if err != nil {
		panic(err)
	}
	settings := cli.New()
	EnvCs.Unlock()

	cfg := new(action.Configuration)
	client := action.NewGet(cfg)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		return nil, err
	}

	res, err := client.Run(name)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func IsExited(name,namespace string) bool {
	res, _ := Get(namespace, name)
	if res != nil {
		return true
	}
	return false
}

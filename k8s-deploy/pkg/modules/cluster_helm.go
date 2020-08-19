package modules

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
)

func (e *Cluster) HelmInstall() error {

	settings, err := service.GetSettings(e.NameSpace)
	if err != nil {
		return err
	}
	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), service.Debug); err != nil {
		return err
	}

	client := action.NewInstall(cfg)

	// if namespace does not exist, create namespace
	//err = CheckNamespace(namespace)
	//if err != nil {
	//	log.Err(err).Msg("CheckNamespace error")
	//	return nil, err
	//}

	// get chart by version from repository
	fc, err := GetFateChart(e.ChartName, e.ChartVersion)
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return err
	}
	log.Debug().Interface("FateChartName", fc.Name).Interface("FateChartVersion", fc.Version).Msg("GetFateChart success")

	// fateChart to helmChart
	chartRequested, err := fc.ToHelmChart()
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return err
	}

	// get values map
	val, err := fc.GetChartValues(e.Spec)
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return err
	}

	log.Debug().Fields(val).Msg("chart values: ")

	client.ReleaseName = e.Name
	client.Namespace = settings.Namespace()

	rel, err := client.Run(chartRequested, val)
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return err
	}
	log.Debug().Interface("runInstall result", rel)

	return nil
}

func (e *Cluster) HelmUpgrade() error {

	settings, err := service.GetSettings(e.NameSpace)
	if err != nil {
		return err
	}

	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), service.Debug); err != nil {
		return err
	}
	client := action.NewUpgrade(cfg)

	client.Namespace = settings.Namespace()

	if client.Version == "" && client.Devel {
		service.Debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	fc, err := GetFateChart(e.ChartName, e.ChartVersion)
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return err
	}
	log.Debug().Interface("FateChartName", fc.Name).Interface("FateChartVersion", fc.Version).Msg("GetFateChart success")

	// fateChart to helmChart
	ch, err := fc.ToHelmChart()
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return err
	}

	// get values map
	val, err := fc.GetChartValues(e.Spec)
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return err
	}
	log.Debug().Fields(val).Msg("chart values: ")

	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return err
		}
	}

	if ch.Metadata.Deprecated {
		fmt.Println("WARNING: This chart is deprecated")
	}

	_, err = client.Run(e.Name, ch, val)
	if err != nil {
		return errors.Wrap(err, "UPGRADE FAILED")
	}
	return nil
}

func (e *Cluster) HelmRollback() error {

	settings, err := service.GetSettings(e.NameSpace)
	if err != nil {
		return err
	}

	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), service.Debug); err != nil {
		return err
	}
	client := action.NewRollback(cfg)
	client.Version = int(e.HelmRevision - 1)
	err = client.Run(e.Name)
	if err != nil {
		return errors.Wrap(err, "UPGRADE FAILED")
	}
	return nil
}

func (e *Cluster) HelmDelete() error {

	settings, err := service.GetSettings(e.NameSpace)
	if err != nil {
		return err
	}

	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), service.Debug); err != nil {
		return err
	}

	client := action.NewUninstall(cfg)

	res, err := client.Run(e.Name)
	if err != nil {
		return err
	}

	log.Debug().Interface("resInfo", res.Info).Msg("delete result")

	return nil
}

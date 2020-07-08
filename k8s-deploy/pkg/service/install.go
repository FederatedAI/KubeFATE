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
package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"os"
	"strconv"
)

type Result struct {
	Namespace    string
	ChartName    string
	ChartVersion string
	ChartValues  map[string]interface{}
	Config       map[string]interface{}
	release      *release.Release
}

// install is create a cluster
// value is a json ,
func Install(namespace, name, chartName, chartVersion string, value *Value) (*Result, error) {

	EnvCs.Lock()
	err := os.Setenv("HELM_NAMESPACE", namespace)
	if err != nil {
		panic(err)
	}
	settings := cli.New()
	EnvCs.Unlock()

	cfg := new(action.Configuration)
	client := action.NewInstall(cfg)

	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		return nil, err
	}

	// if namespace does not exist, create namespace
	//err = CheckNamespace(namespace)
	//if err != nil {
	//	log.Err(err).Msg("CheckNamespace error")
	//	return nil, err
	//}
	// get chart by version from repository
	fc, err := GetFateChart(chartName, chartVersion)
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return nil, err
	}
	log.Debug().Interface("FateChartName", fc.Name).Interface("FateChartVersion", fc.Version).Msg("GetFateChart success")

	// fateChart to helmChart
	chartRequested, err := fc.ToHelmChart()
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return nil, err
	}

	// template to values map
	v, err := value.Unmarshal()
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return nil, err
	}
	log.Debug().Fields(v).Msg("temp values:")

	// get values map
	val, err := fc.GetChartValues(v)
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return nil, err
	}
	// default values
	//val = mergeMaps(chartRequested.Values, val)

	log.Debug().Fields(val).Msg("chart values: ")

	rel, err := runInstall(name, chartRequested, client, val, settings)
	if err != nil {
		log.Err(err).Msg("runInstall error")
		return nil, err
	}

	log.Debug().Interface("runInstall result", rel)

	return &Result{
		Namespace:    settings.Namespace(),
		ChartName:    fc.Name,
		ChartVersion: fc.Version,
		ChartValues:  val,
		release:      rel,
		Config:       v,
	}, nil
}
func newReleaseWriter(releases *release.Release) *releaseElement {
	// Initialize the array so no results returns an empty array instead of null

	r := releases
	element := &releaseElement{
		Name:       r.Name,
		Namespace:  r.Namespace,
		Revision:   strconv.Itoa(r.Version),
		Status:     r.Info.Status.String(),
		Chart:      fmt.Sprintf("%s-%s", r.Chart.Metadata.Name, r.Chart.Metadata.Version),
		AppVersion: r.Chart.Metadata.AppVersion,
	}
	t := "-"
	if tspb := r.Info.LastDeployed; !tspb.IsZero() {
		t = tspb.String()
	}
	element.Updated = t
	return element
}
func runInstall(name string, chartRequested *chart.Chart, client *action.Install, vals map[string]interface{}, settings *cli.EnvSettings) (*release.Release, error) {
	debug("Original chartPath version: %q", client.Version)
	if client.Version == "" && client.Devel {
		debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	client.ReleaseName = name

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		_, _ = fmt.Println("WARNING: This chartPath is deprecated")
	}

	client.Namespace = settings.Namespace()

	return client.Run(chartRequested, vals)
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

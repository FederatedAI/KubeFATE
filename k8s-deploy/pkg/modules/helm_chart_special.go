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

package modules

import (
	"errors"
	"fmt"

	"sigs.k8s.io/yaml"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/chart"
)

func (e *HelmChart) ChartRequestedTohelmChart(chartRequested *chart.Chart) (*HelmChart, error) {
	if chartRequested == nil || chartRequested.Raw == nil {
		log.Error().Msg("chartRequested not exist")
		return nil, errors.New("chartRequested not exist")
	}

	var chartData string
	var valuesData string
	var ValuesTemplate string
	for _, v := range chartRequested.Raw {
		if v.Name == "Chart.yaml" {
			chartData = string(v.Data)
		}
		if v.Name == "values.yaml" {
			valuesData = string(v.Data)
		}
		if v.Name == "values-template.yaml" {
			ValuesTemplate = string(v.Data)
		}
	}

	helmChart := NewHelmChart(chartRequested.Name(),
		chartData, valuesData, chartRequested.Templates, chartRequested.Metadata.Version, chartRequested.AppVersion())

	helmChart.ValuesTemplate = ValuesTemplate
	return helmChart, nil
}

func HelmChartDownload(chartName, chartVersion string) (*HelmChart, error) {
	err := service.RepoAddAndUpdate()
	if err != nil {
		log.Warn().Err(err).Msg("RepoAddAndUpdate error, check kubefate.yaml at env FATECLOUD_REPO_URL values,")
		return nil, err
	}

	chartPath := service.GetChartPath(chartName)
	settings := cli.New()

	cfg := new(action.Configuration)
	client := action.NewInstall(cfg)

	client.ChartPathOptions.Version = chartVersion

	cp, err := client.ChartPathOptions.LocateChart(chartPath, settings)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("FateChart chartPath:", cp).Msg("chartPath:")

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	helmChart, err := ChartRequestedTohelmChart(chartRequested)
	if err != nil {
		return nil, err
	}

	return helmChart, nil
}

func ChartRequestedTohelmChart(chartRequested *chart.Chart) (*HelmChart, error) {
	if chartRequested == nil || chartRequested.Raw == nil {
		log.Error().Msg("chartRequested not exist")
		return nil, errors.New("chartRequested not exist")
	}

	var chartData string
	var valuesData string
	var ValuesTemplate string
	for _, v := range chartRequested.Raw {
		if v.Name == "Chart.yaml" {
			chartData = string(v.Data)
		}
		if v.Name == "values.yaml" {
			valuesData = string(v.Data)
		}
		if v.Name == "values-template.yaml" {
			ValuesTemplate = string(v.Data)
		}
	}

	helmChart := NewHelmChart(chartRequested.Name(),
		chartData, valuesData, chartRequested.Templates, chartRequested.Metadata.Version, chartRequested.AppVersion())

	helmChart.ValuesTemplate = ValuesTemplate
	return helmChart, nil
}

func GetFateChart(chartName, chartVersion string) (*HelmChart, error) {
	hc := &HelmChart{
		Name: chartName, Version: chartVersion,
	}
	if hc.IsExisted() {
		fc, err := hc.Get()
		if err != nil {
			return nil, err
		}
		return &fc, nil
	}
	log.Warn().Str("chartName", chartName).Str("chartVersion", chartVersion).Msg("Chart does not exist in the database")

	helmChart, err := HelmChartDownload(chartName, chartVersion)
	if err != nil {
		return nil, err
	}
	_, err = helmChart.Insert()
	if err != nil {
		return nil, err
	}
	return helmChart, nil
}

func (e *HelmChart) GetChartValuesTemplates() (string, error) {
	if e.ValuesTemplate == "" {
		return "", errors.New("FateChart ValuesTemplate not exist")
	}
	return e.ValuesTemplate, nil
}

func (e *HelmChart) GetChartValues(v map[string]interface{}) (map[string]interface{}, error) {
	// template to values
	template, err := e.GetChartValuesTemplates()
	if err != nil {
		log.Err(err).Msg("GetChartValuesTemplates error")
		return nil, err
	}
	values, err := service.MapToConfig(v, template)
	if err != nil {
		log.Err(err).Interface("v", v).Interface("template", template).Msg("MapToConfig error")
		return nil, err
	}
	// values to map
	vals := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(values), &vals)
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return nil, err
	}
	return vals, nil
}

// todo  get chart by version from repository

func (e *HelmChart) ToHelmChart() (*chart.Chart, error) {
	if e == nil {
		return nil, errors.New("FateChart not exist")
	}
	return e.ConvertToChart()
}

func (e *HelmChart) ConvertToChart() (*chart.Chart, error) {
	c := new(chart.Chart)

	templates := e.Templates

	// Chart file
	chartData := []byte(e.Chart)
	c.Raw = append(c.Raw, &chart.File{Name: "Chart.yaml", Data: chartData})
	if c.Metadata == nil {
		c.Metadata = new(chart.Metadata)
	}
	if err := yaml.Unmarshal(chartData, c.Metadata); err != nil {
		return c, fmt.Errorf("cannot load Chart.yaml %s", err)
	}
	if c.Metadata.APIVersion == "" {
		c.Metadata.APIVersion = chart.APIVersionV1
	}

	// Values file
	valuesData := []byte(e.Values)
	c.Raw = append(c.Raw, &chart.File{Name: "values.yaml", Data: chartData})
	c.Values = make(map[string]interface{})
	if err := yaml.Unmarshal(valuesData, &c.Values); err != nil {
		return c, fmt.Errorf("cannot load Chart.yaml %s", err)
	}

	// Template files
	for _, template := range templates {
		c.Raw = append(c.Raw, template)
		c.Templates = append(c.Templates, &chart.File{Name: template.Name, Data: template.Data})
	}

	// TODO: Handling Chart.lock, values.schema.json, requirements.yaml, requirements.lock files

	return c, nil
}

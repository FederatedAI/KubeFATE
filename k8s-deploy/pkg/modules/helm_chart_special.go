/*
 *  Copyright 2019-2020 VMware, Inc.
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
 */

package modules

import (
	"errors"
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

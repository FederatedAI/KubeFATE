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
	"helm.sh/helm/v3/pkg/cli"
	"os"
)

func Upgrade(namespace, name, chartName, chartVersion string, value *Value) (*Result, error) {

	EnvCs.Lock()
	err := os.Setenv("HELM_NAMESPACE", namespace)
	if err != nil {
		panic(err)
	}
	settings := cli.New()
	EnvCs.Unlock()

	cfg := new(action.Configuration)
	client := action.NewUpgrade(cfg)

	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		return nil, err
	}

	client.Namespace = settings.Namespace()

	if client.Version == "" && client.Devel {
		debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	fc, err := GetFateChart(chartName, chartVersion)
	if err != nil {
		log.Err(err).Msg("GetFateChart error")
		return nil, err
	}
	log.Debug().Interface("FateChartName", fc.Name).Interface("FateChartVersion", fc.Version).Msg("GetFateChart success")

	// fateChart to helmChart
	ch, err := fc.ToHelmChart()
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
	log.Debug().Fields(val).Msg("chart values: ")

	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, err
		}
	}

	if ch.Metadata.Deprecated {
		fmt.Println("WARNING: This chart is deprecated")
	}

	rel, err := client.Run(name, ch, val)
	if err != nil {
		return nil, errors.Wrap(err, "UPGRADE FAILED")
	}
	return &Result{
		Namespace:    settings.Namespace(),
		ChartName:    fc.Name,
		ChartVersion: fc.Version,
		ChartValues:  val,
		release:      rel,
	}, nil
}

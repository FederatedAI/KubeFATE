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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/db"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"

	// "os"
	"testing"

	"github.com/spf13/viper"
	// "helm.sh/helm/v3/pkg/action"
	// "helm.sh/helm/v3/pkg/cli"
)

func TestSaveChartFromPath(t *testing.T) {
	InitConfigForTest()
	// http://github/chart
	// download -> tempPath
	// 1.2  1.3  1.4
	path := "../../fate/"
	helm, err := SaveChartFromPath(path, "fate")
	if err == nil {
		helmUUID, error := db.Save(helm)
		if error == nil {
			t.Log("uuid: ", helmUUID)
		}
	}

}

func TestConvertToChart(t *testing.T) {
	InitConfigForTest()
	helm := &db.HelmChart{}
	result := helm.FindHelmByNameAndVersion("fate", "1.2.0")
	chart, _ := ConvertToChart(result)
	t.Log(chart.AppVersion())

	// settings := cli.New()
	// cfg := new(action.Configuration)
	// client := action.NewInstall(cfg)
	// if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
	// 	return
	// }
	// rel, _ := RunInstall("fate-10000", chart, client, chart.Values, settings)
	// t.Log(rel.Name)
}

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
	logging.InitLog()
}

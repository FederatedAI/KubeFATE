package service

import (
	"fate-cloud-agent/pkg/db"
	"fate-cloud-agent/pkg/utils/config"
	"fate-cloud-agent/pkg/utils/logging"

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

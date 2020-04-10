package service

import (
	"io/ioutil"
	"k8s.io/client-go/rest"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

func MapToConfig(m map[string]interface{}, templates string) (string, error) {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("fate-values-templates").Funcs(funcMap()).Option("missingkey=zero").Parse(string(templates)))

	// Execute the template for each recipient.

	var buf strings.Builder
	err := t.Execute(&buf, m)
	if err != nil {
		log.Error().Msg("executing template:" + err.Error())
		return "", err
	}
	s := strings.ReplaceAll(buf.String(), "<no value>", "")
	return s, nil

}

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	return f
}

func InitKubeConfig() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("~/.kube/config", []byte(config.String()), os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

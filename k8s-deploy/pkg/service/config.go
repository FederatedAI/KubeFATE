package service

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
	"sigs.k8s.io/yaml"
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
func Conf() error {

	values, err := ioutil.ReadFile("E:/machenlong/AI/gitlab/fate-cloud-agent/values.yaml")
	if err != nil {
		log.Error().Msg(err.Error())
	}

	config, err := ioutil.ReadFile("E:/machenlong/AI/gitlab/fate-cloud-agent/config.yaml")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(config, &m)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("values_1.2.0").Funcs(funcMap()).Option("missingkey=zero").Parse(string(values)))
	// Execute the template for each recipient.
	for _, r := range m["PartyList"].([]interface{}) {
		filename := fmt.Sprintf("E:/machenlong/AI/gitlab/fate-cloud-agent/%d-values.yaml", r.(map[interface{}]interface{})["PartyId"])
		file, err := os.Create(filename)
		writer := bufio.NewWriter(file)

		f := make(map[interface{}]interface{})
		for k, v := range m {
			f[k] = v
		}
		for k, v := range r.(map[interface{}]interface{}) {
			f[k] = v
		}
		var buf strings.Builder
		err = t.Execute(&buf, f)
		if err != nil {
			log.Error().Msg("executing template:" + err.Error())
		}
		s := strings.ReplaceAll(buf.String(), "<no value>", "")
		_, _ = writer.WriteString(s)
		err = writer.Flush()
		if err != nil {
			log.Error().Msg("executing template:" + err.Error())
		}
	}
	return nil
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

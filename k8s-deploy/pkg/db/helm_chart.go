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
package db

import (
	"errors"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"helm.sh/helm/v3/pkg/chart"
)

// HelmChart helm chart model
type HelmChart struct {
	Uuid           string        `json:"uuid"`
	Name           string        `json:"name"`
	Chart          string        `json:"chart"`
	Values         string        `json:"values"`
	ValuesTemplate string        `json:"values_template"`
	Templates      []*chart.File `json:"templates"`
	Version        string        `json:"version"`
	AppVersion     string        `json:"app_version"`
}

// NewHelmChart create a new helm chart
func NewHelmChart(name string, chart string, values string, templates []*chart.File, version, appVersion string) *HelmChart {
	helm := &HelmChart{

		Uuid:       uuid.NewV4().String(),
		Name:       name,
		Chart:      chart,
		Values:     values,
		Templates:  templates,
		Version:    version,
		AppVersion: appVersion,
	}

	return helm
}

func (helm *HelmChart) getCollection() string {
	return "helm"
}

// GetUuid get helm uuid
func (helm *HelmChart) GetUuid() string {
	return helm.Uuid
}

// FromBson convert bson to helm
func (helm *HelmChart) FromBson(m *bson.M) (interface{}, error) {
	bsonBytes, err := bson.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(bsonBytes, helm)
	if err != nil {
		return nil, err
	}
	return *helm, nil
}

// FindHelmByNameAndVersion find helm chart via name and version
func (helm *HelmChart) FindHelmByNameAndVersion(name string, version string) *HelmChart {
	filter := bson.M{"name": name, "version": version}
	helms, err := FindByFilter(helm, filter)
	if err == nil && len(helms) != 0 {
		helm0 := helms[0]
		helm0o := helm0.(HelmChart)
		return &helm0o
		// return helms[0].(*HelmChart)
	}
	return nil
}

// FindHelmByNameAndVersion find helm chart via name and version
func FindHelmByNameAndVersion(name, version string) (*HelmChart, error) {
	filter := bson.M{"name": name, "version": version}
	helm := new(HelmChart)
	result, err := FindByFilter(helm, filter)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("not find HelmChart")
	}
	HelmChart := result[0].(HelmChart)
	return &HelmChart, nil

}

// FindHelmChartList find helm chart list
func FindHelmChartList() ([]*HelmChart, error) {
	filter := bson.M{}
	helm := new(HelmChart)
	result, err := FindByFilter(helm, filter)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("not find HelmChart")
	}
	var HelmChartList []*HelmChart
	for _, v := range result {
		r := v.(HelmChart)
		HelmChartList = append(HelmChartList, &r)
	}
	return HelmChartList, nil
}

// FindHelmChartList find helm chart list
func FindHelmChart(chartId string) (*HelmChart, error) {
	filter := bson.M{"uuid": chartId}
	helm := new(HelmChart)
	result, err := FindOneByFilter(helm, filter)
	if err != nil {
		return nil, err
	}

	helmChart, ok := result.(HelmChart)
	if !ok {
		return nil, errors.New("assertion type error")
	}
	return &helmChart, nil
}

func ChartSave(helmChart *HelmChart) (string, error) {

	//find chart
	filter := bson.M{"version": helmChart.Version}
	helmChartFind, err := FindOneByFilter(new(HelmChart), filter)
	if err != nil {
		return "", err
	}
	if helmChartFind == nil {
		helmUUID, err := Save(helmChart)
		if err != nil {
			return "", err
		}

		return helmUUID, err
	}

	err = UpdateByUUID(helmChart, helmChartFind.(HelmChart).Uuid)
	if err != nil {
		return "", err
	}
	return helmChart.Uuid, nil

}

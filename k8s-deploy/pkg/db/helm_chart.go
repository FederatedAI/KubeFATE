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
}

// NewHelmChart create a new helm chart
func NewHelmChart(name string, chart string, values string, templates []*chart.File, version string) *HelmChart {
	helm := &HelmChart{

		Uuid:      uuid.NewV4().String(),
		Name:      name,
		Chart:     chart,
		Values:    values,
		Templates: templates,
		Version:   version,
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
func FindHelmByVersion(version string) (*HelmChart, error) {
	filter := bson.M{"version": version}
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

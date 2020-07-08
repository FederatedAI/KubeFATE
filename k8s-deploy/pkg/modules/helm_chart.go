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
	"database/sql/driver"
	"encoding/json"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"helm.sh/helm/v3/pkg/chart"
)

// HelmChart helm chart model
type HelmChart struct {
	Uuid           string    `json:"uuid" gorm:"type:varchar(36);index:uuid;unique"`
	Name           string    `json:"name" gorm:"type:varchar(16);not null"`
	Chart          string    `json:"chart" gorm:"type:text;not null"`
	Values         string    `json:"values" gorm:"type:text;not null"`
	ValuesTemplate string    `json:"values_template" gorm:"type:text;not null"`
	Templates      Templates `json:"templates" gorm:"type:blob"`
	Version        string    `json:"version" gorm:"type:varchar(32);not null"`
	AppVersion     string    `json:"app_version" gorm:"type:varchar(32);not null"`

	gorm.Model
}
type Templates []*chart.File

type HelmCharts []HelmChart

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

func (s Templates) Value() (driver.Value, error) {
	bJson, err := json.Marshal(s)
	return bJson, err
}

func (s *Templates) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), s)
}

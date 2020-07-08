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
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

type Cluster struct {
	Uuid      string `json:"uuid" gorm:"type:varchar(36);index:uuid;unique"`
	Name      string `json:"name" gorm:"type:varchar(255);not null"`
	NameSpace string `json:"namespaces" gorm:"type:varchar(255);not null"`
	// Cluster revision
	Revision int8 `json:"revision" gorm:"size:8"`
	// Helm chart version, example: fate v1.2.0
	ChartVersion string             `json:"chart_version" gorm:"type:varchar(255);not null"`
	ChartValues  MapStringInterface `json:"chart_values" gorm:"type:blob"`
	// The value of this cluster for installing helm chart
	Values    string             `json:"values" gorm:"type:text"`
	ChartName string             `json:"chart_name" gorm:"type:varchar(255)"`
	Type      string             `json:"cluster_type,omitempty"`
	Metadata  MapStringInterface `json:"metadata" gorm:"type:blob"`
	Status    ClusterStatus      `json:"status"  gorm:"size:8"`
	//Backend          ComputingBackend       `json:"backend"`
	//BootstrapParties Party                  `json:"bootstrap_parties"`
	Config MapStringInterface `json:"Config,omitempty" gorm:"type:blob"`
	Info   MapStringInterface `json:"Info,omitempty" gorm:"type:blob"`

	gorm.Model
}

type MapStringInterface map[string]interface{}

type Clusters []Cluster

type ClusterStatus int8

const (
	Creating_c ClusterStatus = iota + 1
	Deleting_c
	Updating_c
	Running_c
	Unavailable_c
	Deleted_c
)

func (s ClusterStatus) String() string {
	names := map[ClusterStatus]string{
		Creating_c:    "Creating",
		Deleting_c:    "Deleting",
		Updating_c:    "Updating",
		Running_c:     "Running",
		Unavailable_c: "Unavailable",
		Deleted_c:     "Deleted",
	}

	return names[s]
}

// MarshalJSON convert cluster status to string
func (s *ClusterStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (s *ClusterStatus) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	var ClusterStatus ClusterStatus
	switch string(data) {
	case "\"Creating\"":
		ClusterStatus = Creating_c
	case "\"Deleting\"":
		ClusterStatus = Deleting_c
	case "\"Updating\"":
		ClusterStatus = Updating_c
	case "\"Running\"":
		ClusterStatus = Running_c
	case "\"Unavailable\"":
		ClusterStatus = Unavailable_c
	case "\"Deleted\"":
		ClusterStatus = Deleted_c
	default:
		return errors.New("data can't UnmarshalJSON")
	}

	//log.Debug().Interface("JobStatus", JobStatus).Bytes("datab", data).Str("data", string(data)).Msg("UnmarshalJSON")
	*s = ClusterStatus
	return nil
}

// NewCluster create cluster object with basic argument
func NewCluster(name string, nameSpaces, chartName, chartVersion string) *Cluster {
	cluster := &Cluster{
		Uuid:         uuid.NewV4().String(),
		Name:         name,
		NameSpace:    nameSpaces,
		Revision:     0,
		Status:       Creating_c,
		ChartName:    chartName,
		ChartVersion: chartVersion,
	}

	return cluster
}

func (s MapStringInterface) Value() (driver.Value, error) {
	bJson, err := json.Marshal(s)
	return bJson, err
}

func (s *MapStringInterface) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), s)
}

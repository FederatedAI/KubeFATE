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
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"time"
)

type Job struct {
	Uuid      string        `json:"uuid" gorm:"type:varchar(36);index:uuid;unique"`
	StartTime time.Time     `json:"start_time" gorm:"default:Null"`
	EndTime   time.Time     `json:"end_time" gorm:"default:Null"`
	Method    string        `json:"method" gorm:"type:varchar(16);not null"`
	Result    string        `json:"result"  gorm:"type:text"`
	ClusterId string        `json:"cluster_id" gorm:"type:varchar(36)"`
	Creator   string        `json:"creator" gorm:"type:varchar(16);not null"`
	SubJobs   SubJobs       `json:"sub_jobs" gorm:"type:blob"`
	Status    JobStatus     `json:"status"  gorm:"size:8"`
	TimeLimit time.Duration `json:"time_limit"`

	gorm.Model
}

type ClusterArgs struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	ChartName    string `json:"chart_name"`
	ChartVersion string `json:"chart_version"`
	Cover        bool   `json:"cover"`
	Data         []byte `json:"data"`
}

type SubJobs map[string]SubJob

type SubJob struct {
	ModuleName       string
	Status           string
	ModulesPodStatus string
	StartTime        time.Time
	EndTime          time.Time
}

type Jobs []Job

type Method string

const (
	MethodClusterInstall string = "ClusterInstall"
	UNINSTALL
	UPGRADE
	EXEC
)

type JobStatus int8

const (
	JobStatusPending JobStatus = iota + 1
	JobStatusRunning
	JobStatusSuccess
	JobStatusFailed
	JobStatusRollback
	JobStatusTimeout
	JobStatusCanceled
)

func (s JobStatus) String() string {
	names := map[JobStatus]string{
		JobStatusPending:  "Pending",
		JobStatusRunning:  "Running",
		JobStatusSuccess:  "Success",
		JobStatusFailed:   "Failed",
		JobStatusTimeout:  "Timeout",
		JobStatusCanceled: "Canceled",
		JobStatusRollback: "Rollback",
	}

	return names[s]
}

func (s JobStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (s *JobStatus) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	var JobStatus JobStatus
	switch string(data) {
	case "\"Pending\"":
		JobStatus = JobStatusPending
	case "\"Running\"":
		JobStatus = JobStatusRunning
	case "\"Success\"":
		JobStatus = JobStatusSuccess
	case "\"Failed\"":
		JobStatus = JobStatusFailed
	case "\"Timeout\"":
		JobStatus = JobStatusTimeout
	case "\"Canceled\"":
		JobStatus = JobStatusCanceled
	case "\"Rollback\"":
		JobStatus = JobStatusRollback
	default:
		return errors.New("data can't UnmarshalJSON")
	}

	//log.Debug().Interface("JobStatus", JobStatus).Bytes("datab", data).Str("data", string(data)).Msg("UnmarshalJSON")
	*s = JobStatus
	return nil
}

func NewJob(method string, creator string, clusterUuid string) *Job {

	job := &Job{
		Uuid:      uuid.NewV4().String(),
		Method:    method,
		Creator:   creator,
		ClusterId: clusterUuid,
		StartTime: time.Now(),
		Status:    JobStatusPending,
		TimeLimit: 1 * time.Hour,
	}

	return job
}

func (s SubJobs) Value() (driver.Value, error) {
	bJson, err := json.Marshal(s)
	return bJson, err
}

func (s *SubJobs) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), s)
}

func (e *Job) TimeOut() bool {
	return time.Now().After(e.StartTime.Add(e.TimeLimit))
}

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
	uuid "github.com/satori/go.uuid"

	"time"
)

type Job struct {
	Uuid      string        `json:"uuid"`
	StartTime time.Time     `json:"start_time" gorm:"default:Null"`
	EndTime   time.Time     `json:"end_time" gorm:"default:Null"`
	Method    string        `json:"method"`
	Result    string        `json:"result"`
	ClusterId string        `json:"cluster_id"`
	Creator   string        `json:"creator"`
	SubJobs   SubJobs       `json:"sub-jobs" gorm:"type:blob"`
	Status    JobStatus     `json:"status"`
	TimeLimit time.Duration `json:"time_limit"`

	gorm.Model
}

type SubJobs []string

type Jobs []Job

type Method uint32

const (
	INSTALL Method = 1 + iota
	UNINSTALL
	UPGRADE
	EXEC
)

type JobStatus int

const (
	Pending_j JobStatus = iota + 1
	Running_j
	Success_j
	Failed_j
	Retry_j
	Timeout_j
	Canceled_j
)

func (s JobStatus) String() string {
	names := map[JobStatus]string{
		Pending_j:  "Pending",
		Running_j:  "Running",
		Success_j:  "Success",
		Failed_j:   "Failed",
		Retry_j:    "Retry",
		Timeout_j:  "Timeout",
		Canceled_j: "Canceled",
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
		JobStatus = Pending_j
	case "\"Running\"":
		JobStatus = Running_j
	case "\"Success\"":
		JobStatus = Success_j
	case "\"Failed\"":
		JobStatus = Failed_j
	case "\"Retry\"":
		JobStatus = Retry_j
	case "\"Timeout\"":
		JobStatus = Timeout_j
	case "\"Canceled\"":
		JobStatus = Canceled_j
	default:
		return errors.New("data can't UnmarshalJSON")
	}

	//log.Debug().Interface("JobStatus", JobStatus).Bytes("datab", data).Str("data", string(data)).Msg("UnmarshalJSON")
	*s = JobStatus
	return nil
}

func NewJob(method string, creator string) *Job {

	job := &Job{
		Uuid:      uuid.NewV4().String(),
		Method:    method,
		Creator:   creator,
		StartTime: time.Now(),
		Status:    Pending_j,
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
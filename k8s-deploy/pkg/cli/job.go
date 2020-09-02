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
package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/gosuri/uitable"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
	"helm.sh/helm/v3/pkg/cli/output"
)

type Job struct {
}

func (c *Job) getRequestPath() (Path string) {
	return "job/"
}

func (c *Job) addArgs() (Args string) {
	return ""
}

type JobResultList struct {
	Data modules.Jobs
	Msg  string
}

type JobResult struct {
	Data *modules.Job
	Msg  string
}

type JobResultMsg struct {
	Msg string
}

type JobResultErr struct {
	Error string
}

func (c *Job) getResult(Type int) (result interface{}, err error) {
	switch Type {
	case LIST:
		result = new(JobResultList)
	case INFO:
		result = new(JobResult)
	case MSG, JOB:
		result = new(JobResultMsg)
	case ERROR:
		result = new(JobResultErr)
	default:
		err = fmt.Errorf("no type %d", Type)
	}
	return
}

func (c *Job) output(result interface{}, Type int) error {
	switch Type {
	case LIST:
		return c.outPutList(result)
	case INFO:
		return c.outPutInfo(result)
	case MSG, JOB:
		return c.outPutMsg(result)
	case ERROR:
		return c.outPutErr(result)
	default:
		return fmt.Errorf("no type %d", Type)
	}
}

func (c *Job) outPutList(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*JobResultList)
	if !ok {
		return errors.New("type jobResultList not ok")
	}

	joblist := item.Data

	table := uitable.New()
	table.AddRow("UUID", "CREATOR", "METHOD", "STATUS", "STARTTIME", "CLUSTERID", "AGE")
	for _, r := range joblist {
		table.AddRow(r.Uuid, r.Creator, r.Method, r.Status.String(), r.StartTime.Format("2006-01-02 15:04:05"), r.ClusterId, GetDuration(r.StartTime, r.EndTime))
	}
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

func (c *Job) outPutMsg(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*JobResultMsg)
	if !ok {
		return errors.New("type JobResultMsg not ok")
	}

	_, err := fmt.Println(item.Msg)

	return err
}

func (c *Job) outPutErr(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*JobResultErr)
	if !ok {
		return errors.New("type jobResultErr not ok")
	}

	_, err := fmt.Println(item.Error)

	return err
}

func (c *Job) outPutInfo(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	log.Debug().Interface("result", result).Msg("result info")
	item, ok := result.(*JobResult)
	if !ok {
		return errors.New("type jobResult not ok")
	}

	job := item.Data
	log.Debug().Interface("job", job).Msg("job info")

	var subJobs []string
	for _, v := range job.SubJobs {
		subJobs = append(subJobs, fmt.Sprintf("%-20s PodStatus: %s, SubJobStatus: %s, Duration: %6s, StartTime: %s, EndTime: %s",
			v.ModuleName, v.ModulesPodStatus, v.Status, GetDuration(v.StartTime, v.EndTime), v.StartTime.Format("2006-01-02 15:04:05"), v.EndTime.Format("2006-01-02 15:04:05")))
	}

	table := uitable.New()

	colWidth, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	if colWidth == 0 {
		table.MaxColWidth = TableMaxColWidthDefault
	} else {
		table.MaxColWidth = uint(colWidth - ColWidthOffset)
	}
	log.Debug().Int("colWidth", colWidth).Uint("MaxColWidth", table.MaxColWidth).Msg("colWidth ")
	table.Wrap = true // wrap columns
	table.AddRow("UUID", job.Uuid)
	table.AddRow("StartTime", job.StartTime.Format("2006-01-02 15:04:05"))
	table.AddRow("EndTime", job.EndTime.Format("2006-01-02 15:04:05"))
	table.AddRow("Duration", GetDuration(job.StartTime, job.EndTime))
	table.AddRow("Status", job.Status.String())
	table.AddRow("Creator", job.Creator)
	table.AddRow("ClusterId", job.ClusterId)
	table.AddRow("Result", job.Result)
	for i, v := range subJobs {
		if i == 0 {
			table.AddRow("SubJobs", v)
		} else {
			table.AddRow("", v)
		}
	}
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

func GetDuration(startTime, endTime time.Time) string {
	if endTime.IsZero() {
		return HumanDuration(time.Since(startTime))
	}
	return HumanDuration(endTime.Sub(startTime))
}

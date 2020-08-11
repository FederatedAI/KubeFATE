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

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/gosuri/uitable"
	"helm.sh/helm/v3/pkg/cli/output"
)

type Cluster struct {
	all bool
}

func (c *Cluster) getRequestPath() (Path string) {
	return "cluster/"
}
func (c *Cluster) addArgs() (Args string) {

	if c.all {
		Args += "all=true&"
	}

	if len(Args) > 0 {
		Args = "?" + Args
	}
	return Args
}

type ClusterResultList struct {
	Data []*modules.Cluster
	Msg  string
}
type ClusterJobResult struct {
	Data *modules.Job
	Msg  string
}
type ClusterResult struct {
	Data *modules.Cluster
	Msg  string
}

type ClusterResultMsg struct {
	Msg string
}

type ClusterResultErr struct {
	Error string
}

func (c *Cluster) getResult(Type int) (result interface{}, err error) {
	switch Type {
	case LIST:
		result = new(ClusterResultList)
	case INFO:
		result = new(ClusterResult)
	case MSG:
		result = new(ClusterResultMsg)
	case ERROR:
		result = new(ClusterResultErr)
	case JOB:
		result = new(ClusterJobResult)
	default:
		err = fmt.Errorf("no type %d", Type)
	}
	return
}

func (c *Cluster) output(result interface{}, Type int) error {
	switch Type {
	case LIST:
		return c.outPutList(result)
	case INFO:
		return c.outPutInfo(result)
	case MSG:
		return c.outPutMsg(result)
	case ERROR:
		return c.outPutErr(result)
	case JOB:
		return c.outPutJob(result)
	default:
		return fmt.Errorf("no type %d", Type)
	}
}

func (c *Cluster) outPutList(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ClusterResultList)
	if !ok {
		return errors.New("type ClusterResultList not ok")
	}
	table := uitable.New()
	table.AddRow("UUID", "NAME", "NAMESPACE", "REVISION", "STATUS", "CHART", "ChartVERSION", "AGE")
	for _, r := range item.Data {
		table.AddRow(r.Uuid, r.Name, r.NameSpace, r.Revision, r.Status, r.ChartName, r.ChartVersion, HumanDuration(time.Since(r.CreatedAt)))
	}
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

func (c *Cluster) outPutMsg(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ClusterResultMsg)
	if !ok {
		return errors.New("type ClusterResultMsg not ok")
	}

	fmt.Println(item.Msg)
	return nil
}

func (c *Cluster) outPutErr(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ClusterResultErr)
	if !ok {
		return errors.New("type ClusterResultErr not ok")
	}

	fmt.Println(item.Error)

	return nil
}

func (c *Cluster) outPutInfo(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}

	item, ok := result.(*ClusterResult)
	if !ok {
		return errors.New("type ClusterResult not ok")
	}

	cluster := item.Data

	table := uitable.New()
	colWidth, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	if colWidth == 0 {
		table.MaxColWidth = TableMaxColWidthDefault
	} else {
		table.MaxColWidth = uint(colWidth - ColWidthOffset)
	}
	log.Debug().Int("colWidth", colWidth).Uint("MaxColWidth", table.MaxColWidth).Msg("colWidth ")
	table.Wrap = true // wrap columns
	table.AddRow("UUID", cluster.Uuid)
	table.AddRow("Name", cluster.Name)
	table.AddRow("NameSpace", cluster.NameSpace)
	table.AddRow("ChartName", cluster.ChartName)
	table.AddRow("ChartVersion", cluster.ChartVersion)
	table.AddRow("Revision", cluster.Revision)
	table.AddRow("Age", HumanDuration(time.Since(cluster.CreatedAt)))
	table.AddRow("Status", cluster.Status)
	table.AddRow("Values", cluster.Values)
	table.AddRow("Spec", cluster.Spec)
	table.AddRow("Info", cluster.Info)
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

func (c *Cluster) outPutJob(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}

	item, ok := result.(*ClusterJobResult)
	if !ok {
		return errors.New("type ClusterResult not ok")
	}
	fmt.Printf("create job success, job id=%s\r\n", item.Data.Uuid)
	return nil
}

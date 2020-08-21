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

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/gosuri/uitable"
	"helm.sh/helm/v3/pkg/cli/output"
)

type Chart struct {
}

func (c *Chart) getRequestPath() (Path string) {
	return "chart/"
}
func (c *Chart) addArgs() (Args string) {
	return Args
}

type ChartResultList struct {
	Data []*modules.HelmChart
	Msg  string
}

type ChartResult struct {
	Data *modules.HelmChart
	Msg  string
}

type ChartResultMsg struct {
	Msg string
}

type ChartResultErr struct {
	Error string
}

func (c *Chart) getResult(Type int) (result interface{}, err error) {
	switch Type {
	case LIST:
		result = new(ChartResultList)
	case INFO:
		result = new(ChartResult)
	case MSG, JOB:
		result = new(ChartResultMsg)
	case ERROR:
		result = new(ChartResultErr)
	default:
		err = fmt.Errorf("no type %d", Type)
	}
	return
}

func (c *Chart) output(result interface{}, Type int) error {
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
		return fmt.Errorf("output error: no type %d", Type)
	}
}

func (c *Chart) outPutList(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ChartResultList)
	if !ok {
		return errors.New("type ChartResultList not ok")
	}
	table := uitable.New()
	table.AddRow("UUID", "NAME", "VERSION", "APPVERSION")
	for _, r := range item.Data {
		table.AddRow(r.Uuid, r.Name, r.Version, r.AppVersion)
	}
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

func (c *Chart) outPutMsg(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ChartResultMsg)
	if !ok {
		return errors.New("type ChartResultMsg not ok")
	}

	fmt.Println(item.Msg)
	return nil
}

func (c *Chart) outPutErr(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ChartResultErr)
	if !ok {
		return errors.New("type ChartResultErr not ok")
	}

	_, err := fmt.Println(item.Error)

	return err
}

func (c *Chart) outPutInfo(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}

	item, ok := result.(*ChartResult)
	if !ok {
		return errors.New("type ChartResult not ok")
	}

	Chart := item.Data

	table := uitable.New()

	table.AddRow("UUID", Chart.Uuid)
	table.AddRow("Name", Chart.Name)
	table.AddRow("Version", Chart.Version)
	table.AddRow("AppVersion", Chart.AppVersion)
	table.AddRow("Chart", Chart.Chart)
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

package cli

import (
	"errors"
	"fate-cloud-agent/pkg/db"
	"fmt"
	"github.com/gosuri/uitable"
	"helm.sh/helm/v3/pkg/cli/output"
	"os"
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
	Data []*db.HelmChart
	Msg  string
}

type ChartResult struct {
	Data *db.HelmChart
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

func (c *Chart) outPut(result interface{}, Type int) error {
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
		return fmt.Errorf("outPut error: no type %d", Type)
	}
	return nil
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
	table.AddRow("UUID", "NAME", "REVISION")
	for _, r := range item.Data {
		table.AddRow(r.Uuid, r.Name, r.Version)
	}
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

	_, err := fmt.Fprintf(os.Stdout, "%s", item.Msg)

	return err
}

func (c *Chart) outPutErr(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*ChartResultErr)
	if !ok {
		return errors.New("type ChartResultErr not ok")
	}

	_, err := fmt.Fprintf(os.Stdout, "%s", item.Error)

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
	table.AddRow("Chart", Chart.Chart)

	return output.EncodeTable(os.Stdout, table)
}

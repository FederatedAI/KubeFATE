package cli

import (
	"errors"
	"fate-cloud-agent/pkg/db"
	"fmt"
	"github.com/gosuri/uitable"
	"helm.sh/helm/v3/pkg/cli/output"
	"os"
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
	Data []*db.Cluster
	Msg  string
}
type ClusterJobResult struct {
	Data *db.Job
	Msg  string
}
type ClusterResult struct {
	Data *db.Cluster
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

func (c *Cluster) outPut(result interface{}, Type int) error {
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
	return nil
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
	table.AddRow("UUID", "NAME", "NAMESPACE", "REVISION", "STATUS", "CHART", "ChartVERSION")
	for _, r := range item.Data {
		table.AddRow(r.Uuid, r.Name, r.NameSpace, r.Version, r.Status, r.ChartName, r.ChartVersion)
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

	table.AddRow("UUID", cluster.Uuid)
	table.AddRow("Name", cluster.Name)
	table.AddRow("NameSpace", cluster.NameSpace)
	table.AddRow("ChartName", cluster.ChartName)
	table.AddRow("ChartVersion", cluster.ChartVersion)
	table.AddRow("Revision", cluster.Version)
	//table.AddRow("Type", cluster.Type)
	table.AddRow("Status", cluster.Status)
	table.AddRow("Values", cluster.Values)
	table.AddRow("Config", cluster.Config)
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

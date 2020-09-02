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

type Namespace struct {
}

func (c *Namespace) getRequestPath() (Path string) {
	return "namespace/"
}

func (c *Namespace) addArgs() (Args string) {
	return ""
}

type NamespaceResultList struct {
	Data modules.Namespaces
	Msg  string
}

type NamespaceResultErr struct {
	Error string
}

func (c *Namespace) getResult(Type int) (result interface{}, err error) {
	switch Type {
	case LIST:
		result = new(NamespaceResultList)
	case ERROR:
		result = new(NamespaceResultErr)
	default:
		err = fmt.Errorf("no type %d", Type)
	}
	return
}

func (c *Namespace) output(result interface{}, Type int) error {
	switch Type {
	case LIST:
		return c.outPutList(result)
	case ERROR:
		return c.outPutErr(result)
	default:
		return fmt.Errorf("no type %d", Type)
	}
}

func (c *Namespace) outPutList(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*NamespaceResultList)
	if !ok {
		return errors.New("type jobResultList not ok")
	}

	namespacelist := item.Data

	table := uitable.New()
	table.AddRow("NAME", "STATUS", "AGE")
	for _, r := range namespacelist {
		table.AddRow(r.Name, r.Status, r.Age)
	}
	table.AddRow("")
	return output.EncodeTable(os.Stdout, table)
}

func (c *Namespace) outPutErr(result interface{}) error {
	if result == nil {
		return errors.New("no out put data")
	}
	item, ok := result.(*NamespaceResultErr)
	if !ok {
		return errors.New("type NamespaceResultErr not ok")
	}

	_, err := fmt.Println(item.Error)

	return err
}

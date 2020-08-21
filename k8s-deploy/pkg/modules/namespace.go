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
	"fmt"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"k8s.io/apimachinery/pkg/util/duration"
)

type Namespace struct {
	Name        string
	Status      string
	Labels      map[string]string
	Annotations map[string]string
	Age         string
}

type Namespaces []Namespace

func (e *Namespace) GetList() (Namespaces, error) {
	namespaceList, err := service.GetNamespaces()
	if err != nil {
		return nil, err
	}
	var namespaces Namespaces
	for _, v := range namespaceList {
		namespaces = append(namespaces, Namespace{
			Name:        v.Name,
			Status:      fmt.Sprint(v.Status.Phase),
			Labels:      v.Labels,
			Annotations: v.Annotations,
			Age:         duration.HumanDuration(time.Since(v.CreationTimestamp.Time)),
		})
	}
	return namespaces, nil
}

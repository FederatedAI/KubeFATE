/*
 * Copyright 2019-2021 VMware, Inc.
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

package service

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// GetProxySvcNodePorts return rollsite svc NodePort
func GetProxySvcNodePorts(name, namespace string) (int32, error) {
	labelSelector := fmt.Sprintf("name=%s", name)
	svcs, err := KubeClient.GetServices(namespace, labelSelector)
	if err != nil {
		return 0, err
	}

	//svcs.Items[0].GetName()
	for _, v := range svcs.Items {
		if v.GetName() == "rollsite" {
			for _, vv := range v.Spec.Ports {
				if vv.Port == 9370 {
					return vv.NodePort, nil
				}
			}
		}
	}
	return 0, nil
}

// GetServiceStatus func
func GetServiceStatus(Services *v1.ServiceList) map[string]string {
	status := make(map[string]string)
	for _, v := range Services.Items {
		status[v.Name] = v.Status.String()
	}
	return status
}

// GetClusterServiceStatus func
func GetClusterServiceStatus(name, namespace string) (map[string]string, error) {
	labelSelector := fmt.Sprintf("name=%s", name)
	list, err := KubeClient.GetServices(namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	return GetServiceStatus(list), nil
}

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

	v1 "k8s.io/api/apps/v1"
)

// GetStsList gets the statefulSets list under the namespace
func GetStsList(clusterName, namespace string) (*v1.StatefulSetList, error) {

	list, err := KubeClient.GetStsList(namespace, getLabelSelector(namespace, clusterName))
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetStsStatus gets the status if a certain statefulSet
func GetStsStatus(sts *v1.StatefulSet) (string, string) {
	if sts.Status.ReadyReplicas >= sts.Status.Replicas {
		return "Available", "all the replicas are in the ready state"
	} else {
		return "Progressing", "Detailed status need to be checked by kubectl CLI"
	}
}

// GetStssStatus gets the status of a list of statefulSets
func GetStssStatus(stss *v1.StatefulSetList) (map[string]string, error) {
	status := make(map[string]string)
	for _, v := range stss.Items {
		Type, _ := GetStsStatus(&v)
		status[v.Name] = fmt.Sprintf("%s", Type)
	}
	return status, nil
}

// GetClusterStsStatus gets all the statefulSet related information with the cluster name and namespace
func GetClusterStsStatus(name, namespace string) (map[string]string, error) {
	stsList, err := GetStsList(name, namespace)
	if err != nil {
		return nil, err
	}
	return GetStssStatus(stsList)
}

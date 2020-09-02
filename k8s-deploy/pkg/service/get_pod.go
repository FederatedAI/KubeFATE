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
package service

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(namespace, LabelSelector string) (*v1.PodList, error) {
	clientset, err := getClientset()
	if err != nil {
		fmt.Println(err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: LabelSelector})
	return pods, err
}

func GetClusterStatus(name, namespace string) (map[string]string, error) {
	var labelSelector string
	labelSelector = fmt.Sprintf("name=%s", name)
	list, err := GetPods(namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	return GetPodStatus(list), nil
}
func GetPodStatus(pods *v1.PodList) map[string]string {

	status := make(map[string]string)
	for _, v := range pods.Items {
		for _, vv := range v.Status.ContainerStatuses {
			if vv.State.Running != nil {
				status[vv.Name] = "Running"
				continue
			}
			if vv.State.Waiting != nil {
				status[vv.Name] = "Waiting"
				continue
			}
			if vv.State.Terminated != nil {
				status[vv.Name] = "Terminated"
				continue
			}
			status[vv.Name] = "Unknown"
		}
	}
	return status

}

func CheckClusterStatus(ClusterStatus map[string]string) bool {
	if len(ClusterStatus) == 0 {
		return false
	}
	var clusterStatusOk = true
	for _, v := range ClusterStatus {
		if v != "Running" {
			clusterStatusOk = false
		}
	}
	return clusterStatusOk
}

func GetPodList(name, namespace string) ([]string, error) {
	var labelSelector string
	labelSelector = fmt.Sprintf("name=%s", name)
	list, err := GetPods(namespace, labelSelector)
	if err != nil {
		return nil, err
	}
	var podList []string
	for _, v := range list.Items {
		podList = append(podList, v.GetName())
	}
	return podList, nil
}

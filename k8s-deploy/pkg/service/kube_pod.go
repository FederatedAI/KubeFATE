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
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// GetClusterPodStatus GetClusterPodStatus
func GetClusterPodStatus(name, namespace string) (map[string]string, error) {

	list, err := KubeClient.GetPods(namespace, getLabelSelector(namespace, name))
	if err != nil {
		return nil, err
	}

	return GetPodStatus(list), nil
}

//GetPodStatus GetPodStatus
func GetPodStatus(pods *corev1.PodList) map[string]string {

	status := make(map[string]string)
	for _, v := range pods.Items {
		switch string(v.Status.Phase) {
		case "Running", "Succeeded", "Pending", "Failed":
			status[v.Name] = string(v.Status.Phase)
			continue
		default:
			status[v.Name] = "Unknown"
		}
	}
	return status
}

// GetPodList GetPodList
func GetPodList(name, namespace string) ([]string, error) {

	list, err := KubeClient.GetPods(namespace, getLabelSelector(namespace, name))
	if err != nil {
		return nil, err
	}
	var podList []string
	for _, v := range list.Items {
		podList = append(podList, v.GetName())
	}
	return podList, nil
}

// GetPodNameByModule is Get Pod By Module
func GetPodNameByModule(namespace, name, modules string) (string, error) {
	labelSelector := getLabelSelector(namespace, name)
	podList, err := KubeClient.GetPods(namespace, labelSelector)
	if err != nil {
		return "", err
	}

	for _, pod := range podList.Items {
		for _, container := range pod.Spec.Containers {
			if container.Name == modules {
				return pod.Name, nil
			}
		}
	}

	return "", errors.New("module no find")
}

// getLabelSelector is Get LabelSelector
// This part depends on matchLabels of helm hart _helpers.tpl file
func getLabelSelector(namespace, name string) string {

	return fmt.Sprintf("name=%s", name)
}

// getPodContainerList getPodContainerList
// return map[ContainerName]podName
func getPodContainerList(name, namespace, container string) (map[string]string, error) {

	list, err := KubeClient.GetPods(namespace, getLabelSelector(namespace, name))
	if err != nil {
		return nil, err
	}
	var podContainerList = make(map[string]string)
	for _, v := range list.Items {
		for _, vv := range v.Spec.Containers {
			if container == "" {
				podContainerList[vv.Name] = v.GetName()
			} else {
				if container == vv.Name {
					podContainerList[vv.Name] = v.GetName()
				}
			}
		}

	}
	return podContainerList, nil
}

//GetPodContainersStatus GetPodContainersStatus
func GetPodContainersStatus(ClusterName, namespace string) (map[string]string, error) {
	list, err := KubeClient.GetPods(namespace, getLabelSelector(namespace, ClusterName))
	if err != nil {
		return nil, err
	}
	var podContainerList = make(map[string]string)
	for _, v := range list.Items {
		for _, vv := range v.Spec.Containers {
			podContainerList[vv.Name] = fmt.Sprintf("%s", v.Status.Phase)
		}
	}
	return podContainerList, nil
}

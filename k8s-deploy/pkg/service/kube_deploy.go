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

	corev1 "k8s.io/api/core/v1"
)

// GetDeployList GetDeployList
func GetDeployList(clusterName, namespace string) (*v1.DeploymentList, error) {

	list, err := KubeClient.GetDeploymentList(namespace, getLabelSelector(namespace, clusterName))
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetDeployStatus GetDeployStatus
func GetDeployStatus(deploy *v1.Deployment) (string, string) {

	for _, v := range deploy.Status.Conditions {
		if v.Type == v1.DeploymentAvailable && v.Status == corev1.ConditionTrue {
			return fmt.Sprint(v1.DeploymentAvailable), v.Message
		}
	}
	for _, v := range deploy.Status.Conditions {
		if v.Type == v1.DeploymentProgressing && v.Status == corev1.ConditionTrue {
			return fmt.Sprint(v1.DeploymentProgressing), v.Message
		}
	}
	for _, v := range deploy.Status.Conditions {
		if v.Type == v1.DeploymentReplicaFailure && v.Status == corev1.ConditionTrue {
			return fmt.Sprint(v1.DeploymentReplicaFailure), v.Message
		}
	}
	return "Undefined", fmt.Sprintf("please use kubectl cli check deploy status of %s", deploy.Name)
}

func GetDeploymentStatus(deploys *v1.DeploymentList) (map[string]string, error) {
	status := make(map[string]string)
	for _, v := range deploys.Items {
		Type, _ := GetDeployStatus(&v)
		status[v.Name] = fmt.Sprintf("%s", Type)
	}
	return status, nil
}

// GetClusterDeployStatus GetClusterDeployStatus
func GetClusterDeployStatus(name, namespace string) (map[string]string, error) {
	deploymentList, err := GetDeployList(name, namespace)
	if err != nil {
		return nil, err
	}
	return GetDeploymentStatus(deploymentList)
}

func CheckStatus(status string) bool {
	return status == "Available"
}

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

package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment Deployment
type Deployment interface {
	GetDeployment(namespace, deploymentName string) (*v1.Deployment, error)
	GetDeploymentList(namespace, LabelSelector string) (*v1.DeploymentList, error)
}

// GetDeployment is get a Deployment
func (e *Kube) GetDeployment(namespace, deploymentName string) (*v1.Deployment, error) {
	deployment, err := e.client.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	return deployment, err
}

// GetDeploymentList is get a GetDeploymentList
func (e *Kube) GetDeploymentList(namespace, LabelSelector string) (*v1.DeploymentList, error) {
	deploymentlist, err := e.client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: LabelSelector})
	return deploymentlist, err
}

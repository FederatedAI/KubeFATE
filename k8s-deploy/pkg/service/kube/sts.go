/*
 * Copyright 2019-2022 VMware, Inc.
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

// Sts StatefulSet
type Sts interface {
	GetSts(namespace, stsName string) (*v1.StatefulSet, error)
	GetStsList(namespace, LabelSelector string) (*v1.StatefulSetList, error)
}

// GetSts gets a StatefulSet
func (e *Kube) GetSts(namespace, stsName string) (*v1.StatefulSet, error) {
	sts, err := e.client.AppsV1().StatefulSets(namespace).Get(context.Background(), stsName, metav1.GetOptions{})
	return sts, err
}

// GetStsList gets a StatefulSet lis
func (e *Kube) GetStsList(namespace, LabelSelector string) (*v1.StatefulSetList, error) {
	stsList, err := e.client.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: LabelSelector})
	return stsList, err
}

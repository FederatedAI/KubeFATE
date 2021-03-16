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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Services interface
type Services interface {
	GetServices(namespace, labelSelector string) (*corev1.ServiceList, error)
}

// GetServices is get Services list
func (e *Kube) GetServices(namespace, labelSelector string) (*corev1.ServiceList, error) {
	return e.client.CoreV1().Services(namespace).List(e.ctx, metav1.ListOptions{LabelSelector: labelSelector})
}

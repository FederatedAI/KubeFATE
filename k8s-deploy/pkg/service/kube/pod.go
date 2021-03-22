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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod interface
type Pod interface {
	GetPod(podName, namespace string) (*corev1.Pod, error)
	GetPods(namespace, LabelSelector string) (*corev1.PodList, error)
}

// GetPod is get a pod info
func (e *Kube) GetPod(podName, namespace string) (*corev1.Pod, error) {
	pod, err := e.client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	return pod, err
}

// GetPods is get pod list info
func (e *Kube) GetPods(namespace, LabelSelector string) (*corev1.PodList, error) {
	pods, err := e.client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: LabelSelector})
	return pods, err
}

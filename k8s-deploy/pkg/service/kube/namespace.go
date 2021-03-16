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

// Namespace interface
type Namespace interface {
	GetNamespace(namespace string) (*corev1.Namespace, error)
	GetNamespaces() (*corev1.NamespaceList, error)
	CreateNamespace(namespaceName string) (*corev1.Namespace, error)
}

// GetNamespace is get a pod info
func (e *Kube) GetNamespace(namespace string) (*corev1.Namespace, error) {
	return e.client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
}

// GetNamespaces is get a pod info
func (e *Kube) GetNamespaces() (*corev1.NamespaceList, error) {
	return e.client.CoreV1().Namespaces().List(e.ctx, metav1.ListOptions{})
}

// CreateNamespace Create a Namespace
func (e *Kube) CreateNamespace(namespaceName string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	return e.client.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
}

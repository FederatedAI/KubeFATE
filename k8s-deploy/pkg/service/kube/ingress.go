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

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ingress interface
type Ingress interface {
	GetIngress(ingressName, namespace string) (*networkingv1.Ingress, error)
	GetIngresses(namespace, labelSelector string) (*networkingv1.IngressList, error)
}

// GetIngress is get a Ingress
func (e *Kube) GetIngress(ingressName, namespace string) (*networkingv1.Ingress, error) {
	ingress, err := e.client.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	return ingress, err
}

// GetIngresses is get list of Ingress
func (e *Kube) GetIngresses(namespace, labelSelector string) (*networkingv1.IngressList, error) {
	return e.client.NetworkingV1().Ingresses(namespace).List(e.ctx, metav1.ListOptions{LabelSelector: labelSelector})
}

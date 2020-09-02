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

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngress(name, namespace string) (*v1beta1.IngressList, error) {
	clientset, err := getClientset()
	if err != nil {
		fmt.Println(err)
	}

	ingressList, err := clientset.ExtensionsV1beta1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{})

	return ingressList, err
}

func GetIngressUrl(name, namespace string) ([]string, error) {
	var urls []string

	ingressList, err := GetIngress(name, namespace)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, ingress := range ingressList.Items {
		for _, v := range ingress.Spec.Rules {
			urls = append(urls, v.Host)
		}
	}

	return urls, nil
}

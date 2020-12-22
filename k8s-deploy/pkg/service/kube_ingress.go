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
	"fmt"
)

// GetIngressURLList is Get Ingress Url list
func GetIngressURLList(name, namespace string) ([]string, error) {
	var urls []string
	labelSelector := getLabelSelector(namespace, name)
	ingressList, err := KubeClient.GetIngresses(namespace, labelSelector)
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

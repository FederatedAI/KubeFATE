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
	v1 "k8s.io/api/core/v1"
)

// GetNamespaces GetNamespaces
func GetNamespaces() ([]v1.Namespace, error) {
	namespaceList, err := KubeClient.GetNamespaces()
	if err != nil {
		return nil, err
	}
	return namespaceList.Items, nil
}

// CreateNamespace CreateNamespace
func CreateNamespace(namespace string) error {
	_, err := KubeClient.CreateNamespace(namespace)
	return err
}

// CheckNamespace CheckNamespace
func CheckNamespace(namespace string) error {

	namespaces, err := KubeClient.GetNamespace(namespace)
	if err == nil {
		return nil
	}

	if namespaces != nil {
		// namespace exist
		return nil
	}

	_, err = KubeClient.CreateNamespace(namespace)

	return err

}

// getDefaultNamespace Get Default Namespace
func getDefaultNamespace(namespace string) string {
	if namespace != "" {
		return namespace
	}

	// TODO: Get the default namespace
	return "default"
}

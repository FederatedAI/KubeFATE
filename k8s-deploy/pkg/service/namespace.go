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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespace(namespace string) (string, error) {
	clientset, err := getClientset()

	if err != nil {
		return "", err
	}
	rNamespace, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})

	if err != nil {
		return "", err
	}
	fmt.Printf("rNamespace :%+v\n", rNamespace)
	return rNamespace.Name, nil
}

func GetNamespaceList() ([]string, error) {
	clientset, err := getClientset()

	if err != nil {
		return nil, err
	}

	ncl, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []string
	for _, v := range ncl.Items {
		namespaces = append(namespaces, v.Name)
	}

	return namespaces, nil
}

func GetNamespaces() ([]v1.Namespace, error) {
	clientset, err := getClientset()

	if err != nil {
		return nil, err
	}

	ncl, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ncl.Items, nil
}

func CreateNamespace(namespace string) error {
	clientset, err := getClientset()

	if err != nil {
		return err
	}
	nc := new(v1.Namespace)

	nc.TypeMeta = metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"}
	nc.ObjectMeta = metav1.ObjectMeta{
		Name: namespace,
	}
	nc.Spec = v1.NamespaceSpec{}
	_, err = clientset.CoreV1().Namespaces().Create(context.Background(), nc, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	return nil
}

func CheckNamespace(namespace string) error {
	clientset, err := getClientset()
	if err != nil {
		return err
	}
	namespaces, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err == nil {
		return nil
	}

	if namespaces != nil {
		// namespace exist
		return nil
	}

	nc := new(v1.Namespace)

	nc.TypeMeta = metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"}
	nc.ObjectMeta = metav1.ObjectMeta{
		Name: namespace,
	}
	nc.Spec = v1.NamespaceSpec{}
	_, err = clientset.CoreV1().Namespaces().Create(context.Background(), nc, metav1.CreateOptions{})

	return err

}

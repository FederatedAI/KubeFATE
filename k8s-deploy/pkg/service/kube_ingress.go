package service

import (
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngress(name, namespace string) (*v1beta1.IngressList, error) {
	clientset, err := getClientset()
	if err != nil {
		fmt.Println(err)
	}

	ingressList, err := clientset.ExtensionsV1beta1().Ingresses(namespace).List(metav1.ListOptions{})

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

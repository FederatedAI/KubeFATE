package service

import (
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngress(name, namespace string) (*v1beta1.Ingress, error) {
	clientset, err := getClientset()
	if err != nil {
		fmt.Println(err)
	}

	ingressName := "fateboard"

	ingress, err := clientset.ExtensionsV1beta1().Ingresses(namespace).Get(ingressName, metav1.GetOptions{})
	return ingress, err
}

func GetIngressUrl(name, namespace string) (string,error) {

	ingress, err := GetIngress(name, namespace)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return ingress.Spec.Rules[0].Host, nil
}

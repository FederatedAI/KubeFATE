package service

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespace(namespace string) (string, error) {
	clientset, err := getClientset()

	if err != nil {
		return "", err
	}
	rNamespace, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})

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

	ncl, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []string
	for _, v := range ncl.Items {
		namespaces = append(namespaces, v.Name)
	}

	return namespaces, nil
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
	_, err = clientset.CoreV1().Namespaces().Create(nc)

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
	namespaces, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
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
	_, err = clientset.CoreV1().Namespaces().Create(nc)

	return err

}

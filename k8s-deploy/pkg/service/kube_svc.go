package service

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetProxySvcNodePorts(name, namespace string) ([]int32, error) {
	var labelSelector string
	labelSelector = fmt.Sprintf("name=%s", name)
	svcs, err := GetServices(namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	var nodePorts []int32

	//svcs.Items[0].GetName()
	for _, v := range svcs.Items {
		if v.GetName() == "proxy" {
			for _, vv := range v.Spec.Ports {
				nodePorts = append(nodePorts, vv.NodePort)
			}
		}
	}
	return nodePorts, nil
}

func GetServices(namespace, LabelSelector string) (*v1.ServiceList, error) {
	clientset, err := getClientset()
	if err != nil {
		return nil, err
	}

	svcs, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector: LabelSelector})
	return svcs, err
}

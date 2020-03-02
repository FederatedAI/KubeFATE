package service

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeIp() ([]string, error) {
	svcs, err := GetNodes()
	if err != nil {
		return nil, err
	}

	var nodeIp []string
	for _, v := range svcs.Items {
		nodeIp = append(nodeIp,v.Status.Addresses[0].Address)
	}
	return nodeIp, nil
}

func GetNodes() (*v1.NodeList, error) {
	clientset, err := getClientset()
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	return nodes, err
}

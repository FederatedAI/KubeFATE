package service

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(namespace string) (*v1.PodList, error) {
	clientset, err := getClientset()
	if err != nil {
		fmt.Println(err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	return  pods, err

}

func checkPodStatus(pods *v1.PodList, ) bool {
	for _, v := range pods.Items {
		if v.Status.Phase != v1.PodRunning {
			return false
		}
	}
	return true

}

// todo get pod by name
func CheckClusterStatus(name, namespace string) (bool, error) {
	list, err := GetPods(namespace)
	if err != nil {
		return false, err
	}

	return checkPodStatus(list), nil
}

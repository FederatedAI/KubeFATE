package service

import (
	"github.com/spf13/viper"
	"sync"

	"helm.sh/helm/v3/pkg/kube"
	"k8s.io/client-go/kubernetes"
)

var EnvCs sync.Mutex

func getClientset() (*kubernetes.Clientset, error) {
	configFlags := kube.GetConfig(viper.GetString("kube.config"), viper.GetString("kube.context"), viper.GetString("kube.namespace"))
	config, _ := configFlags.ToRESTConfig()
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

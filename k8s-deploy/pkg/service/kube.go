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
	"os"
	"sync"

	"helm.sh/helm/v3/pkg/cli"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service/kube"
)

var EnvCs sync.Mutex

type kubeClient interface {
	kube.Pod
	kube.Namespace
	kube.Ingress
	kube.Node
	kube.Services
	kube.Log
	kube.Deployment
}

var KubeClient kubeClient = &kube.KUBE

func GetSettings(namespace string) (*cli.EnvSettings, error) {
	EnvCs.Lock()
	err := os.Setenv("HELM_NAMESPACE", namespace)
	if err != nil {
		return nil, err
	}
	settings := cli.New()
	err = os.Unsetenv("HELM_NAMESPACE")
	if err != nil {
		return nil, err
	}
	EnvCs.Unlock()

	return settings, nil
}

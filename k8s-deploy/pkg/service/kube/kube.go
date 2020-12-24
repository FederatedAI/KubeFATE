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

package kube

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/kube"
	"k8s.io/client-go/kubernetes"
)

// Kube struct
type Kube struct {
	client kubernetes.Interface
	ctx    context.Context
}

// KUBE Kube
var KUBE Kube

func getClientset() (*kubernetes.Clientset, error) {
	configFlags := kube.GetConfig(viper.GetString("kube.config"), viper.GetString("kube.context"), viper.GetString("kube.namespace"))
	config, _ := configFlags.ToRESTConfig()
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

func init() {
	client, err := getClientset()
	if err != nil {
		log.Error().Err(err).Msg("getClientset")
	}
	KUBE.client = client
	KUBE.ctx = context.Background()
}

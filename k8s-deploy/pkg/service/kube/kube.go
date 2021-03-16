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

package kube

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
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

	config, err := getConfig(viper.GetString("kube.namespace"), viper.GetString("kube.context"), viper.GetString("kube.config")).ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

func getConfig(namespace, context, kubeConfig string) *genericclioptions.ConfigFlags {
	cf := genericclioptions.NewConfigFlags(true)
	cf.Namespace = &namespace
	cf.Context = &context
	cf.KubeConfig = &kubeConfig
	return cf
}

func init() {
	client, err := getClientset()
	if err != nil {
		log.Error().Err(err).Msg("getClientset")
	}
	KUBE.client = client
	KUBE.ctx = context.Background()
}

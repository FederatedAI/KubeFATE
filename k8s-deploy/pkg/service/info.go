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

import "github.com/rs/zerolog/log"

// GetClusterInfo GetClusterInfo
func GetClusterInfo(name, namespace string) (map[string]interface{}, error) {
	ip, err := GetNodeIP()
	if err != nil {
		log.Error().Str("func", "GetNodeIP()").Err(err).Msg("GetNodeIP error")
		return nil, err
	}
	port, err := GetProxySvcNodePorts(name, getDefaultNamespace(namespace))
	if err != nil {
		log.Error().Str("func", "GetProxySvcNodePorts()").Err(err).Msg("GetProxySvcNodePorts error")
		return nil, err
	}

	containerList, err := GetPodContainersStatus(name, getDefaultNamespace(namespace))
	if err != nil {
		log.Error().Str("func", "GetPodContainersStatus()").Err(err).Msg("GetPodContainersStatus error")
		return nil, err
	}

	deploymentList, err := GetClusterDeployStatus(name, getDefaultNamespace(namespace))
	if err != nil {
		log.Error().Str("func", "GetClusterDeployStatus()").Err(err).Msg("GetClusterDeployStatus error")
		return nil, err
	}

	status := make(map[string]interface{})

	status["containers"] = containerList
	status["deployments"] = deploymentList

	ingressURLList, err := GetIngressURLList(name, getDefaultNamespace(namespace))
	if err != nil {
		log.Error().Str("func", "GetIngressURLList()").Err(err).Msg("GetIngressURLList error")
		return nil, err
	}

	info := make(map[string]interface{})

	if len(ip) > 0 {
		info["ip"] = ip[len(ip)-1]
	}
	if port != 0 {
		info["port"] = port
	}

	info["status"] = status

	info["dashboard"] = ingressURLList

	log.Debug().Interface("cluster-info", info).Msg("show the cluster info real-time status")

	return info, nil
}

//GetClusterStatus GetClusterStatus
func GetClusterStatus(name, namespace string) (map[string]string, error) {
	return GetClusterDeployStatus(name, namespace)
}

// CheckClusterStatus CheckClusterStatus
func CheckClusterStatus(ClusterStatus map[string]string) bool {
	if len(ClusterStatus) == 0 {
		return false
	}
	var clusterStatusOk = true
	for _, v := range ClusterStatus {
		if !CheckStatus(v) {
			clusterStatusOk = false
		}
	}
	return clusterStatusOk
}

// CheckClusterStatus CheckClusterStatus
func CheckClusterInfoStatus(ClusterInfoStatus map[string]interface{}) bool {
	Status, ok := ClusterInfoStatus["status"]
	if !ok {
		return false
	}

	deployments, ok := Status.(map[string]interface{})["deployments"]
	if !ok {
		return false
	}

	ClusterStatus := deployments.(map[string]string)

	if len(ClusterStatus) == 0 {
		return false
	}
	var clusterStatusOk = true
	for _, v := range ClusterStatus {
		if !CheckStatus(v) {
			clusterStatusOk = false
		}
	}
	return clusterStatusOk
}

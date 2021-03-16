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

// GetClusterInfo GetClusterInfo
func GetClusterInfo(name, namespace string) (map[string]interface{}, error) {
	ip, err := GetNodeIP()
	if err != nil {
		return nil, err
	}
	port, err := GetProxySvcNodePorts(name, getDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}
	podList, err := GetPodList(name, getDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}

	ingressURLList, err := GetIngressURLList(name, getDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}

	info := make(map[string]interface{})

	if len(ip) > 0 {
		info["ip"] = ip[len(ip)-1]
	}
	if port != 0 {
		info["port"] = port
	}
	info["pod"] = podList

	info["dashboard"] = ingressURLList

	return info, nil
}

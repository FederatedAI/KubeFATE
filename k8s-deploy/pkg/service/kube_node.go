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

// GetNodeIP get a node ip
func GetNodeIP() ([]string, error) {

	svcs, err := KubeClient.GetNodes("")
	if err != nil {
		return nil, err
	}

	var nodeIP []string
	for _, v := range svcs.Items {
		nodeIP = append(nodeIP, v.Status.Addresses[0].Address)
	}
	return nodeIP, nil
}

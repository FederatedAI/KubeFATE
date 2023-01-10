/*
 * Copyright 2019-2022 VMware, Inc.
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

package job

import "github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"

type UpgradeManager interface {
	validate(specOld, specNew modules.MapStringInterface) error
	getCluster(specOld, specNew modules.MapStringInterface) modules.Cluster
	waitFinish(interval, round int) bool
}

type FallbackUpgradeManager struct {
	UpgradeManager
}

func (um *FallbackUpgradeManager) validate(specold, specNew modules.MapStringInterface) error {
	return nil
}

func (um *FallbackUpgradeManager) getCluster(specold, specNew modules.MapStringInterface) modules.Cluster {
	return modules.Cluster{
		Name: "fallbackUM",
	}
}

func (um *FallbackUpgradeManager) waitFinish(interval, round int) bool {
	return true
}

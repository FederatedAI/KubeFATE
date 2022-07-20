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

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"strings"
)

func getMysqlCredFromSpec(clusterSpec modules.MapStringInterface) (username, password string) {
	mysqlSpec := clusterSpec["mysql"].(modules.MapStringInterface)
	if mysqlSpec["user"] == nil {
		username = "fate"
	} else {
		username = mysqlSpec["user"].(string)
	}
	if mysqlSpec["password"] == nil {
		password = "fate_dev"
	} else {
		password = mysqlSpec["password"].(string)
	}
	return
}

func constructFumSpec(oldSpec, newSpec modules.MapStringInterface) (fumSpec modules.MapStringInterface) {
	oldVersion := strings.ReplaceAll(oldSpec["chartVersion"].(string), "v", "")
	newVersion := strings.ReplaceAll(newSpec["chartVersion"].(string), "v", "")
	mysqlUsername, mysqlPassword := getMysqlCredFromSpec(newSpec)
	res := modules.MapStringInterface{
		"username": mysqlUsername,
		"password": mysqlPassword,
		"start":    oldVersion,
		"target":   newVersion,
	}
	return res
}

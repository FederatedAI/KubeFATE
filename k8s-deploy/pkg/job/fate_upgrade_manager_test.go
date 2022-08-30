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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	specOld := modules.MapStringInterface{
		"chartName":    "fate",
		"chartVersion": "v1.7.2",
	}
	specNew := modules.MapStringInterface{
		"chartName":    "fate",
		"chartVersion": "v1.8.0",
	}
	fum := FateUpgradeManager{
		namespace: "blabla",
	}
	// Happy path
	err := fum.validate(specOld, specNew)
	assert.Nil(t, err)

	// framework different
	specNew["chartName"] = "openfl"
	err = fum.validate(specOld, specNew)
	assert.NotNil(t, err)
	specNew["chartName"] = "fate"

	// spec identical
	specNew["chartVersion"] = "v1.7.2"
	err = fum.validate(specOld, specNew)
	assert.NotNil(t, err)
	specNew["chartVersion"] = "v1.8.0"

	// do not support downgrade
	specNew["chartVersion"] = "v1.6.0"
	err = fum.validate(specOld, specNew)
	assert.NotNil(t, err)
	specNew["chartVersion"] = "v1.8.0"

	// fate version < 1.7.1
	specOld["chartVersion"] = "v1.7.0"
	err = fum.validate(specOld, specNew)
	assert.NotNil(t, err)
}

func TestGetCluster(t *testing.T) {
	specOld := modules.MapStringInterface{
		"chartName":    "fate",
		"chartVersion": "v1.7.2",
	}
	specNew := modules.MapStringInterface{
		"chartName":    "fate",
		"chartVersion": "v1.8.0",
		"mysql": map[string]interface{}{
			"user":     "fate",
			"password": "fate_dev",
		},
	}
	fum := FateUpgradeManager{
		namespace: "blabla",
	}
	actual := fum.getCluster(specOld, specNew)
	expect := modules.Cluster{
		Name:         "fate-upgrade-manager",
		NameSpace:    "blabla",
		ChartName:    "fate-upgrade-manager",
		ChartVersion: "v1.0.0",
		Spec: modules.MapStringInterface{
			"password": "fate_dev",
			"username": "fate",
			"start":    "1.7.2",
			"target":   "1.8.0",
		},
	}
	assert.True(t, reflect.DeepEqual(actual, expect))
}

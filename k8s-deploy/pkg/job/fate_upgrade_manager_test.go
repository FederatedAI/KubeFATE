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

func TestConstructFumSpec(t *testing.T) {
	oldSpec := modules.MapStringInterface{
		"chartVersion": "v1.7.2",
	}
	newSpec := modules.MapStringInterface{
		"chartVersion": "v1.8.0",
		"mysql": modules.MapStringInterface{
			"user":     "fate",
			"password": "fate_dev",
		},
	}
	actual := ConstructFumSpec(oldSpec, newSpec)
	expect := modules.MapStringInterface{
		"password": "fate_dev",
		"username": "fate",
		"start":    "1.7.2",
		"target":   "1.8.0",
	}
	assert.True(t, reflect.DeepEqual(actual, expect))
}

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

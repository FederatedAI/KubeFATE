package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateVersion(t *testing.T) {
	actual := validateVersion("v1.7.0", "1.7.0-release")
	assert.True(t, actual)

	actual = validateVersion("v1.7.0", "1.7.1-release")
	assert.False(t, actual)

	actual = validateVersion("v1.6.0", "1.6.0-release")
	assert.True(t, actual)
}

func Test_getUpgradeScripts(t *testing.T) {
	actual, err := getUpgradeScripts("v1.7.0", "v1.7.1")
	expected := []string{"v1.7.0-to-v1.7.1.sql"}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)

	actual, err = getUpgradeScripts("v1.7.1", "v1.9.0")
	expected = []string{
		"v1.7.1-to-v1.7.2.sql",
		"v1.7.2-to-v1.8.0.sql",
		"v1.8.0-to-v1.9.0.sql",
	}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)

	actual, err = getUpgradeScripts("v1.8.0", "v1.10.0")
	expected = []string{}
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
}

package job

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateVersion(t *testing.T) {
	err := validateFateVersion("v1.7.0", "1.7.0-release")
	assert.Nil(t, err)

	err = validateFateVersion("v1.7.0", "1.7.1-release")
	assert.NotNil(t, err)

	err = validateFateVersion("v1.6.0", "1.6.0-release")
	assert.Nil(t, err)
}

func Test_getUpgradeScripts(t *testing.T) {
	viper.Set("upgradesupportedfateversions", []string{
		"1.7.0", "1.7.1", "1.7.2", "1.8.0", "1.9.0",
	})
	actual, err := getUpgradeScripts("v1.7.0", "v1.7.1")
	expected := []string{"1.7.0-to-1.7.1.sql"}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)

	actual, err = getUpgradeScripts("v1.7.1", "v1.9.0")
	expected = []string{
		"1.7.1-to-1.7.2.sql",
		"1.7.2-to-1.8.0.sql",
		"1.8.0-to-1.9.0.sql",
	}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)

	actual, err = getUpgradeScripts("v1.8.0", "v1.10.0")
	expected = []string{}
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
}

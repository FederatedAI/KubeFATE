package job

import (
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

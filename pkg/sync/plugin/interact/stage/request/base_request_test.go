package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitPluginMap(t *testing.T) {
	assert.NotEqual(t, nil, Plugins, "failed to init Request plugin map")
}

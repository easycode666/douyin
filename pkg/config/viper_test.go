package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	config := InitConfig("user-config")
	assert.Equal(t, "DouyinUserServer", config.GetString("Server.Name"))
}

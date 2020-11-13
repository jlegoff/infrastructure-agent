package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfigJsonParse(t *testing.T) {
	jsonFile := []byte(`{
"config_version": "1",
"config": {
  "databind": {
    "labels": {
      "integration_group": "my group",
      "environment": "production"
    },
    "discovery": {
      "ttl": "1s"
    },
    "variables": {
      "myVariable": {
        "ttl": "1s"
      }
    }
  },
  "integrations": [
    {
      "exec": "/path/to/executable"
    },
    {
      "exec": "/path/to/another/executable"
    }]
  }
}`)
	config := ConfigProtocol{}
	require.NoError(t, json.Unmarshal(jsonFile, &config))
	assert.Equal(t, "1s", config.Config.Databind.Discovery.TTL)
	require.Contains(t, config.Config.Databind.Variables, "myVariable")
	assert.Equal(t, "1s", config.Config.Databind.Variables["myVariable"].TTL)
	assert.Len(t, config.Config.Integrations, 2)
	assert.Contains(t, config.Config.Integrations, ConfigEntry{Exec: ShlexOpt{"/path/to/executable"}})
	assert.Contains(t, config.Config.Integrations, ConfigEntry{Exec: ShlexOpt{"/path/to/another/executable"}})
}
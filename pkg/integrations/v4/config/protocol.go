package config

import (
	"encoding/json"
)

// json sent to the agent to run integrations
type ConfigProtocol struct {
	ConfigProtocolVersion string `json:"config_version"`
	Config                YAML   `json:"config"`
}

func IsConfigRequest(line []byte) (isConfig bool, configProtocol ConfigProtocol) {
	if err := json.Unmarshal(line, &configProtocol); err != nil {
		return
	}
	isConfig = true
	return
}

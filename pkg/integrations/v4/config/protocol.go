package config

// json sent to the agent to run integrations
type ConfigProtocol struct {
	ConfigProtocolVersion string `json:"config_version"`
	Config                YAML   `json:"config"`
}

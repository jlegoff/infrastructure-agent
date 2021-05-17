// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
package config

import (
	"github.com/kelseyhightower/envconfig"
	"path/filepath"
)

const (
	defaultConnectEnabled = true
)

func init() {
	defaultConfigFiles = []string{
		"newrelic-infra.yml",
		filepath.Join("/etc", "newrelic-infra.yml"),
		filepath.Join("/etc", "newrelic-infra", "newrelic-infra.yml"),
	}
	defaultPluginConfigFiles = []string{
		filepath.Join("/etc", "newrelic-infra-plugins.yml"),
		filepath.Join("/etc", "newrelic-infra", "newrelic-infra-plugins.yml"),
	}
	defaultPluginInstanceDir = filepath.Join("/etc", "newrelic-infra", "integrations.d")
	defaultConfigDir = filepath.Join("/etc", "newrelic-infra")

	defaultAgentDir = filepath.Join("/var", "db", "newrelic-infra")
	defaultLogFile = filepath.Join("/var", "db", "newrelic-infra", "newrelic-infra.log")
	defaultNetworkInterfaceFilters = map[string][]string{
		"prefix":  {"dummy", "lo", "vmnet", "sit", "tun", "tap", "veth"},
		"index-1": {"tun", "tap"},
	}

	defaultLoggingBinDir = "logging"
	defaultLoggingConfigsDir = "logging.d"
	defaultFluentBitExe = "fluent-bit"
	defaultFluentBitParsers = "parsers.conf"
	defaultFluentBitNRLib = "out_newrelic.so"

	// this is the default dir the infra sdk uses to store "temporary" data
	defaultIntegrationsTempDir = filepath.Join("/tmp", "nr-integrations")
}

func runtimeValues() (userMode, agentUser, executablePath string) {
	return ModeRoot, "", ""
}

func configOverride(cfg *Config) {
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		clog.WithError(err).Error("unable to interpret environment variables")
	}
}

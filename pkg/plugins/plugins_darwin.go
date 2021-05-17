// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
package plugins

import (
	"github.com/newrelic/infrastructure-agent/internal/agent"
	pluginsLinux "github.com/newrelic/infrastructure-agent/internal/plugins/linux"
	"github.com/newrelic/infrastructure-agent/pkg/integrations/v4/emitter"
	"github.com/newrelic/infrastructure-agent/pkg/metrics"
	metricsSender "github.com/newrelic/infrastructure-agent/pkg/metrics/sender"
	"github.com/newrelic/infrastructure-agent/pkg/metrics/storage"
	"github.com/newrelic/infrastructure-agent/pkg/plugins/ids"
	"github.com/newrelic/infrastructure-agent/pkg/plugins/proxy"
)

func RegisterPlugins(a *agent.Agent, em emitter.Emitter) error {
	a.RegisterPlugin(NewHostAliasesPlugin(a.Context, a.GetCloudHarvester()))
	config := a.Context.Config()

	if config.ProxyConfigPlugin {
		a.RegisterPlugin(proxy.ConfigPlugin(a.Context))
	}
	a.RegisterPlugin(NewCustomAttrsPlugin(a.Context))
	a.RegisterPlugin(NewAgentConfigPlugin(*ids.NewPluginID("metadata", "agent_config"), a.Context))
	a.DeprecatePlugin(ids.PluginID{"hostinfo", "hostinfo"})
	a.RegisterPlugin(pluginsLinux.NewHostinfoPlugin(a.Context, a.GetCloudHarvester()))
	if config.HTTPServerEnabled {
		httpSrv, err := NewHTTPServerPlugin(a.Context, config.HTTPServerHost, config.HTTPServerPort, em)
		if err != nil {
			slog.
				WithField("port", config.HTTPServerPort).
				WithField("host", config.HTTPServerHost).
				WithError(err).
				Error("cannot create HTTP server")
		} else {
			a.RegisterPlugin(httpSrv)
		}
	}

	if config.FilesConfigOn {
		a.RegisterPlugin(NewConfigFilePlugin(*ids.NewPluginID("files", "config"), a.Context))
	}

	sender := metricsSender.NewSender(a.Context)
	storageSampler := storage.NewSampler(a.Context)
	systemSampler := metrics.NewSystemSampler(a.Context, storageSampler)
	sender.RegisterSampler(systemSampler)
	sender.RegisterSampler(storageSampler)
	a.RegisterMetricsSender(sender)

	return nil
}

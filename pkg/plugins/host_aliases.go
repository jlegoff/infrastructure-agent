// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
package plugins

import (
	"fmt"
	"github.com/newrelic/infrastructure-agent/pkg/entity"
	"time"

	"github.com/newrelic/infrastructure-agent/pkg/sysinfo"
	"github.com/newrelic/infrastructure-agent/pkg/sysinfo/cloud"

	"github.com/newrelic/infrastructure-agent/pkg/plugins/ids"
	"github.com/newrelic/infrastructure-agent/pkg/sysinfo/hostname"

	"github.com/newrelic/infrastructure-agent/internal/agent"
	"github.com/newrelic/infrastructure-agent/pkg/config"
	"github.com/newrelic/infrastructure-agent/pkg/log"
)

type CloudAlias struct {
	source 	string
	alias 	string
}

type HostAliasesPlugin struct {
	agent.PluginCommon
	resolver       hostname.Resolver
	cloudHarvester cloud.Harvester // Used to get metadata for the instance.
	logger         log.Entry
}

func NewHostAliasesPlugin(ctx agent.AgentContext, cloudHarvester cloud.Harvester) agent.Plugin {
	id := ids.PluginID{
		Category: "metadata",
		Term:     "host_aliases",
	}
	return &HostAliasesPlugin{
		PluginCommon: agent.PluginCommon{
			ID:      id,
			Context: ctx},
		cloudHarvester: cloudHarvester,
		resolver:       ctx.HostnameResolver(),
		logger:         slog.WithField("id", id),
	}
}

func (self *HostAliasesPlugin) getHostAliasesDataset() (dataset agent.PluginInventoryDataset, err error) {
	fullHostname, shortHostname, err := self.resolver.Query()
	if err != nil {
		return nil, fmt.Errorf("error resolving hostname: %s", err)
	}
	if len(fullHostname) == 0 {
		return nil, fmt.Errorf("retrieved empty hostname")
	}

	dataset = append(dataset, sysinfo.HostAliases{
		Alias:  fullHostname,
		Source: sysinfo.HOST_SOURCE_HOSTNAME,
	})

	dataset = append(dataset, sysinfo.HostAliases{
		Alias:  shortHostname,
		Source: sysinfo.HOST_SOURCE_HOSTNAME_SHORT,
	})

	// Retrieve the host alias from config
	if self.Context.Config().DisplayName != "" {
		dataset = append(dataset, sysinfo.HostAliases{
			Alias:  self.Context.Config().DisplayName,
			Source: sysinfo.HOST_SOURCE_DISPLAY_NAME,
		})
	}

	// Retrieve the instance ID if the host happens to be running in a cloud VM. If we hit an
	// error do not return the dataset because this would make the agent reconnect under a different id, triggering
	// HostNotReporting alerts. See https://github.com/newrelic/infrastructure-agent/issues/94
	if self.shouldCollectCloudMetadata() {
		cloudAlias, err := self.collectCloudMetadata()
		if err != nil {
			self.logger.WithError(err).Debug("Could not retrieve instance ID. Either this is not the cloud or the metadata API returned an error.")
			return nil, err
		}
		dataset = append(dataset, sysinfo.HostAliases{
			Source: cloudAlias.source,
			Alias:  cloudAlias.alias,
		})
	}

	return dataset, nil
}

// shouldCollectCloudMetadata will check if we should query for the cloud metadata.
func (self *HostAliasesPlugin) shouldCollectCloudMetadata() bool {
	return !self.Context.Config().DisableCloudMetadata &&
		!self.Context.Config().DisableCloudInstanceId &&
		self.cloudHarvester.GetCloudType().ShouldCollect()
}

// Collect cloud metadata and set self.cloudAlias to include whatever we found
func (self *HostAliasesPlugin) collectCloudMetadata() (alias *CloudAlias, err error) {
	instanceID, err := self.cloudHarvester.GetInstanceID()
	if err != nil {
		return nil, err
	}

	cloudAlias := CloudAlias {
		source: self.cloudHarvester.GetCloudSource(),
		alias: 	instanceID,
	}
	return &cloudAlias, nil
}

func (self *HostAliasesPlugin) Run() {
	refreshTimer := time.NewTicker(1)

	for {
		select {
		case <-refreshTimer.C:
			refreshTimer.Stop()
			refreshTimer = time.NewTicker(config.FREQ_PLUGIN_HOST_ALIASES * time.Second)
			{
				var dataset agent.PluginInventoryDataset
				var err error
				self.logger.Debug("Starting harvest.")
				if dataset, err = self.getHostAliasesDataset(); err != nil {
					self.logger.WithError(err).Error("fetching aliases")
					continue
				}
				self.logger.WithField("dataset", dataset).Debug("Completed harvest, emitting.")
				self.EmitInventory(dataset, entity.NewFromNameWithoutID(self.Context.AgentIdentifier()))
				self.logger.Debug("Completed emitting.")
			}
		}
	}
}

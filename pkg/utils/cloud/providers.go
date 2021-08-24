/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cloud

import (
	"github.com/go-logr/logr"

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/RasaHQ/rasactl/pkg/utils/cloud/providers"
)

type Provider struct {
	Name       types.CloudProvider
	Log        logr.Logger
	ExternalIP string
}

func (p *Provider) New() types.CloudProvider {
	provider := p.detect(
		providers.Google(),
		providers.Amazon(),
		providers.Azure(),
		providers.Alibaba(),
		providers.Oracle(),
		providers.DigitalOcean(),
	)

	p.Name = provider

	switch provider {
	case types.CloudProviderGoogle:
		p.ExternalIP = providers.GoogleGetExternalIP()
	case types.CloudProviderAmazon:
		p.ExternalIP = providers.AmazonGetExternalIP()
	case types.CloudProviderAzure:
		p.ExternalIP = providers.AzureGetExternalIP()
	case types.CloudProviderDigitalOcean:
		p.ExternalIP = providers.DigitalOceanGetExternalIP()
	case types.CloudProviderAlibaba:
		p.ExternalIP = providers.AlibabaGetExternalIP()
	default:
		p.ExternalIP = "0.0.0.0"
	}

	p.Log.Info("Detecting cloud provider", "provider", string(provider), "externalIP", p.ExternalIP)

	return provider
}

func (p *Provider) detect(providers ...types.CloudProvider) types.CloudProvider {

	for _, provider := range providers {
		if provider != types.CloudProviderUnknown {
			return provider
		}
	}

	return types.CloudProviderUnknown
}

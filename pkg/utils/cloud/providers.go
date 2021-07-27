package cloud

import (
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils/cloud/providers"
	"github.com/go-logr/logr"
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

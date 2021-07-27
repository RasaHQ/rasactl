package cloud

import (
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils/cloud/providers"
	"github.com/go-logr/logr"
)

type Provider struct {
	Name types.CloudProvider
	Log  logr.Logger
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
	p.Log.Info("Detecting cloud provider", "provider", string(provider))

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

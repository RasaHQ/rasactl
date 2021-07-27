package providers

import (
	"io/ioutil"
	"strings"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Amazon() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/product_version")
	if strings.Contains(string(data), "amazon") {
		return types.CloudProviderAmazon
	}
	return types.CloudProviderUnknown
}

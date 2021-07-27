package providers

import (
	"io/ioutil"
	"strings"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Alibaba() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/product_name")
	if strings.Contains(string(data), "Alibaba Cloud") {
		return types.CloudProviderAmazon
	}
	return types.CloudProviderUnknown
}

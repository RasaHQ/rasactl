package providers

import (
	"io/ioutil"
	"strings"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Google() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/product_name")
	if strings.Contains(string(data), "Google") {
		return types.CloudProviderGoogle
	}
	return types.CloudProviderUnknown
}

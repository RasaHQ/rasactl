package providers

import (
	"io/ioutil"
	"strings"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func DigitalOcean() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if strings.Contains(string(data), "DigitalOcean") {
		return types.CloudProviderDigitalOcean
	}
	return types.CloudProviderUnknown
}

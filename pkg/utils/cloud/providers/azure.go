package providers

import (
	"io/ioutil"
	"strings"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Azure() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if strings.Contains(string(data), "Microsoft Corporation") {
		return types.CloudProviderAzure
	}
	return types.CloudProviderUnknown
}

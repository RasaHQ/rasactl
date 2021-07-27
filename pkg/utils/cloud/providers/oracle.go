package providers

import (
	"io/ioutil"
	"strings"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Oracle() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/chassis_asset_tag")
	if strings.Contains(string(data), "OracleCloud") {
		return types.CloudProviderOracle
	}
	return types.CloudProviderUnknown
}

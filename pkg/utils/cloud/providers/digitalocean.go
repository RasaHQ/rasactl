package providers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func DigitalOcean() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if strings.Contains(string(data), "DigitalOcean") {
		return types.CloudProviderDigitalOcean
	}
	return types.CloudProviderUnknown
}

func DigitalOceanGetExternalIP() string {
	var body []byte
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 20,
	}
	req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address", nil)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		body = bodyData
	}
	return string(body)
}

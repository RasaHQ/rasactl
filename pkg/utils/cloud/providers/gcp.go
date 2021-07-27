package providers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Google() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/product_name")
	if strings.Contains(string(data), "Google") {
		return types.CloudProviderGoogle
	}
	return types.CloudProviderUnknown
}

func GoogleGetExternalIP() string {
	var body []byte
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 20,
	}
	req, _ := http.NewRequest("GET", "http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip", nil)
	req.Header.Add("Metadata-Flavor", "Google")
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

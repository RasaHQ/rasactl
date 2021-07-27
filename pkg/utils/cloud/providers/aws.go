package providers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Amazon() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/product_version")
	if strings.Contains(string(data), "amazon") {
		return types.CloudProviderAmazon
	}
	return types.CloudProviderUnknown
}

func AmazonGetExternalIP() string {
	var body []byte
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 20,
	}
	req, _ := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/public-ipv4", nil)
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

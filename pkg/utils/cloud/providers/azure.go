package providers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/types"
)

func Azure() types.CloudProvider {
	data, _ := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if strings.Contains(string(data), "Microsoft Corporation") {
		return types.CloudProviderAzure
	}
	return types.CloudProviderUnknown
}

func AzureGetExternalIP() string {
	var body []byte
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 20,
	}
	req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance/network/interface/0/ipv4/ipAddress/0/publicIpAddress?api-version=2017-08-01&format=text", nil)
	req.Header.Add("Metadata", "true")
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

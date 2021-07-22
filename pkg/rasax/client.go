package rasax

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

type RasaX struct {
	URL            string
	Token          string
	Log            logr.Logger
	SpinnerMessage *status.SpinnerMessage
	WaitTimeout    time.Duration
	client         *http.Client
}

func (r *RasaX) New() {
	r.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 120,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 20 * time.Second,
			}).Dial,
		},
	}
}

func (r *RasaX) GetHealthEndpoint() (*types.HealthEndpointsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/health", r.URL), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 304 || resp.StatusCode == 502 {
		bodyData := &types.HealthEndpointsResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &bodyData)
		return bodyData, nil

	} else {
		return nil, errors.Errorf("The Rasa X health endpoint has returned status code %s", resp.Status)
	}
}

func (r *RasaX) GetVersionEndpoint() (*types.VersionEndpointResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/version", r.URL), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyData := &types.VersionEndpointResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &bodyData)
		return bodyData, nil

	} else {
		return nil, errors.Errorf("The Rasa X health endpoint has returned status code %s", resp.Status)
	}
}

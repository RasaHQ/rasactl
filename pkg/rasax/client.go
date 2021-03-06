/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package rasax

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"

	"github.com/RasaHQ/rasactl/pkg/status"
	"github.com/RasaHQ/rasactl/pkg/types"
	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

// RasaX defines Rasa X client.
type RasaX struct {
	// BearerToken stores a bearer token.
	BearerToken string

	// Flags defines the command flags.
	Flags *types.RasaCtlFlags

	// URL is a Rasa X URL
	URL string

	// Token stores a Rasa X admin token.
	Token string

	// Log defines logger.
	Log logr.Logger

	// Log defines the spinner object.
	SpinnerMessage *status.SpinnerMessage

	// WaitTimeout defines timeout for the client.
	WaitTimeout time.Duration

	client *http.Client
}

// New initializes a new Rasa X client.
func (r *RasaX) New() {
	r.client = &http.Client{
		Timeout: time.Second * 30,
	}
}

func (r *RasaX) getURL() string {
	urlAddress := r.URL

	if !utils.IsURLAccessible(urlAddress) {
		parsedURL, _ := url.Parse(urlAddress)

		urlAddress = fmt.Sprintf("%s://%s", parsedURL.Scheme, "127.0.0.1")
		if parsedURL.Port() != "" {
			urlAddress = fmt.Sprintf("%s:%s", urlAddress, parsedURL.Port())
		}
		r.Log.V(1).Info("The URL is not accessible for the health endpoint, using internal address", "url", r.URL, "internalURL", urlAddress)
	}

	return urlAddress
}

// GetHealthEndpoint returns a response from the /api/health endpoint.
func (r *RasaX) GetHealthEndpoint() (*rtypes.HealthEndpointsResponse, error) {
	urlAddress := r.getURL()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/health", urlAddress), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 304 || resp.StatusCode == 502 {
		bodyData := &rtypes.HealthEndpointsResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &bodyData); err != nil {
			return nil, err
		}
		return bodyData, nil

	}
	return nil, xerrors.Errorf("The Rasa X health endpoint has returned status code %s", resp.Status)
}

func (r *RasaX) GetVersionEndpoint() (*rtypes.VersionEndpointResponse, error) {
	urlAddress := r.getURL()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/version", urlAddress), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyData := &rtypes.VersionEndpointResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &bodyData); err != nil {
			return nil, err
		}
		return bodyData, nil

	}
	return nil, xerrors.Errorf("The Rasa X health endpoint has returned status code %s", resp.Status)
}

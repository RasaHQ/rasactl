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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
)

func (r *RasaX) Auth(username, password string) (*rtypes.AuthEndpointResponse, error) {
	values := map[string]string{
		"username": username,
		"password": password,
	}

	jsonValue, _ := json.Marshal(values)

	urlAddress := r.getURL()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth", urlAddress), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		bodyData := &rtypes.AuthEndpointResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &bodyData); err != nil {
			return nil, err
		}
		return bodyData, nil

	case 401:
		return nil, errors.Errorf("Unauthorized")

	default:
		return nil, errors.Errorf("The Rasa X health endpoint has returned status code %s", resp.Status)
	}
}

// ValidateToken validates token and returns true if a given token is valid.
func (r *RasaX) ValidateToken(token string) bool {
	urlAddress := r.getURL()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/config", urlAddress), nil)
	if err != nil {
		r.Log.V(1).Info("Can't validate token", "error", err)
		return false
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := r.client.Do(req)
	if err != nil {
		r.Log.V(1).Info("Can't validate token", "error", err)
		return false
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return true
	case 401:
		r.Log.Info("Token is invalid", "token", token)
		return false
	default:
		return false
	}
}

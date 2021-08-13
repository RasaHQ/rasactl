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

	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
	"github.com/pkg/errors"
)

func (r *RasaX) Auth(username, password string) (*rtypes.AuthEndpointResponse, error) {
	urlAddress := r.getURL()
	values := map[string]string{
		"username": username,
		"password": password,
	}

	jsonValue, _ := json.Marshal(values)

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
		json.Unmarshal(body, &bodyData)
		return bodyData, nil

	case 401:
		return nil, errors.Errorf("Unauthorized")

	default:
		return nil, errors.Errorf("The Rasa X health endpoint has returned status code %s", resp.Status)
	}
}

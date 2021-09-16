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
)

// SaveEnvironments add environments to Rasa X via the /environments endpoint.
// Required Rasa X >= 1.0.
func (r *RasaX) SaveEnvironments(body []rtypes.EnvironmentsEndpointRequest) error {
	urlAddress := r.getURL()
	b := new(bytes.Buffer)

	if err := json.NewEncoder(b).Encode(body); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/environments", urlAddress)
	r.Log.V(1).Info("Sending a request to Rasa X", "url", url)
	request, err := http.NewRequest("PUT", url, b)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return nil
	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s", content)
	}
}

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

	"golang.org/x/xerrors"
)

// EnterpriseActivate activates an Enterprise license via the /api/license endpoint.
func (r *RasaX) EnterpriseActivate(license string) error {
	urlAddress := r.getURL()
	url := fmt.Sprintf("%s/api/license", urlAddress)
	r.Log.V(1).Info("Sending a request to Rasa X", "url", url)

	values := map[string]string{
		"license": license,
	}

	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		fmt.Printf("\nThe Enterprise license has been activated.\n")
		return nil

	case 401:
		return xerrors.Errorf("Unauthorized")

	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return xerrors.Errorf("The Rasa X license endpoint has returned status code %s, body: %s", resp.Status, content)
	}
}

// EnterpriseDeactivate deactivates an Enterprise license via the /api/license endpoint.
func (r *RasaX) EnterpriseDeactivate() error {
	urlAddress := r.getURL()
	url := fmt.Sprintf("%s/api/license", urlAddress)
	r.Log.V(1).Info("Sending a request to Rasa X", "url", url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 204:
		fmt.Println("The Enterprise license has been deactivated.")
		return nil

	case 401:
		return xerrors.Errorf("Unauthorized")

	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return xerrors.Errorf("The Rasa X license endpoint has returned status code %s, body: %s", resp.Status, content)
	}
}

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
package rasactl

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/RasaHQ/rasactl/pkg/status"
)

func (r *RasaCtl) checkIfRasaOSSProductionIsConnected() error {
	r.initRasaXClient()

	resp, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	if resp.Rasa.Production == "0.0.0" {
		return fmt.Errorf("rasa server is not connected to the production environment")
	}

	return nil
}

func (r *RasaCtl) ModelUpload() error {

	if err := r.checkIfRasaOSSProductionIsConnected(); err != nil {
		return err
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	err = r.RasaXClient.ModelUpload()
	return err
}

func (r *RasaCtl) ModelDelete() error {
	if err := r.checkIfRasaOSSProductionIsConnected(); err != nil {
		return err
	}

	resp, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	if resp.Rasa.Production == "0.0.0" {
		return fmt.Errorf("rasa server is not connected to the production environment")
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	err = r.RasaXClient.ModelDelete()
	return err
}

func (r *RasaCtl) ModelDownload() error {
	if err := r.checkIfRasaOSSProductionIsConnected(); err != nil {
		return err
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	err = r.RasaXClient.ModelDownload()
	return err
}

func (r *RasaCtl) ModelTag() error {
	if err := r.checkIfRasaOSSProductionIsConnected(); err != nil {
		return err
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	err = r.RasaXClient.ModelTag()
	return err
}

func (r *RasaCtl) ModelList() error {
	data := [][]string{}
	header := []string{"Name", "Version", "Compatible", "Tags", "Hash", "Trained At"}

	if err := r.checkIfRasaOSSProductionIsConnected(); err != nil {
		return err
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	models, err := r.RasaXClient.ModelList()
	if err != nil {
		return err
	}

	if len(models.Models) == 0 {
		fmt.Println("Nothing to show, upload model to see results.")
		return nil
	}

	for _, model := range models.Models {
		sec, dec := math.Modf(model.TrainedAt)
		tags := "none"
		if len(model.Tags) != 0 {
			tags = strings.Join(model.Tags, ",")
		}
		data = append(data, []string{
			model.Model,
			model.Version,
			fmt.Sprintf("%t", model.IsCompatible),
			tags,
			model.Hash,
			time.Unix(int64(sec), int64(dec*(1e9))).Format("02 Jan 06 15:04 MST"),
		})
	}
	status.PrintTable(
		header,
		data,
	)
	return nil
}

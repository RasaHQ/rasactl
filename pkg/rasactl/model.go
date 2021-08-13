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

func (r *RasaCtl) ModelUpload() error {
	r.initRasaXClient()

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	if err := r.RasaXClient.ModelUpload(); err != nil {
		return err
	}
	return nil
}

func (r *RasaCtl) ModelDelete() error {
	r.initRasaXClient()

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	if err := r.RasaXClient.ModelDelete(); err != nil {
		return err
	}
	return nil
}

func (r *RasaCtl) ModelDownload() error {
	r.initRasaXClient()

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	if err := r.RasaXClient.ModelDownload(); err != nil {
		return err
	}
	return nil
}

func (r *RasaCtl) ModelTag() error {
	r.initRasaXClient()

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	if err := r.RasaXClient.ModelTag(); err != nil {
		return err
	}
	return nil
}

func (r *RasaCtl) ModelList() error {
	data := [][]string{}
	header := []string{"Name", "Version", "Compatible", "Tags", "Hash", "Trained At"}

	r.initRasaXClient()

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
		data = append(data, []string{
			model.Model,
			model.Version,
			fmt.Sprintf("%t", model.IsCompatible),
			strings.Join(model.Tags, ","),
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

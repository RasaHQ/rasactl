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

	"golang.org/x/xerrors"

	"github.com/RasaHQ/rasactl/pkg/utils"
)

// EnterpriseActivate activates an Enterprise license.
func (r *RasaCtl) EnterpriseActivate() error {
	r.initRasaXClient()

	version, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	if version.Enterprise {
		fmt.Println("The Enterprise license is already active.")
		return nil
	}

	if utils.RasaXVersionConstrains(version.RasaX, "< 1.0.0") {
		return xerrors.Errorf("this command is available for Rasa X 1.0.0 or newer")
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	license, err := utils.ReadLicense(r.Flags)
	if err != nil {
		return err
	}

	return r.RasaXClient.EnterpriseActivate(license)
}

// EnterpriseDeactivate deactivates an Enterprise license.
func (r *RasaCtl) EnterpriseDeactivate() error {
	r.initRasaXClient()

	version, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	if !version.Enterprise {
		return xerrors.Errorf("an Enterprise license is not active")
	}

	if utils.RasaXVersionConstrains(version.RasaX, "< 1.0.0") {
		return xerrors.Errorf("this command is available for Rasa X 1.0.0 or newer")
	}

	token, err := r.getAuthToken()
	if err != nil {
		return err
	}
	r.RasaXClient.BearerToken = token

	return r.RasaXClient.EnterpriseDeactivate()
}

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

	"github.com/RasaHQ/rasactl/pkg/credentials"
	"github.com/RasaHQ/rasactl/pkg/credentials/helpers"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

func (r *RasaCtl) AuthLogin() error {

	if r.isLogged() {
		fmt.Println("Already logged.")
		return nil
	}

	username, password, err := utils.ReadCredentials(r.Flags)
	if err != nil {
		return err
	}

	r.Log.Info("Initializing Rasa X client")
	r.initRasaXClient()

	r.Log.Info("Getting a token")

	authRes, err := r.RasaXClient.Auth(username, password)
	if err != nil {
		return err
	}

	credsStore := credentials.Credentials{
		Namespace: r.Namespace,
		Helper:    helpers.Helper,
	}

	r.Log.Info("Storing credentials in the store", "name", "rasactl-login", "namespace", r.Namespace)
	if err := credsStore.Set("rasactl-login", username, password); err != nil {
		return err
	}
	r.Log.Info("Storing credentials in the store", "name", "rasactl-token", "namespace", r.Namespace)
	if err := credsStore.Set("rasactl-token", r.Namespace, authRes.AccessToken); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Successfully logged.")

	return nil
}

func (r *RasaCtl) AuthLogout() error {
	r.initRasaXClient()

	credsStore := credentials.Credentials{
		Namespace: r.Namespace,
		Helper:    helpers.Helper,
	}

	r.Log.Info("Deleting credentials from the store", "name", "rasactl-login", "namespace", r.Namespace)
	if err := credsStore.Delete("rasactl-login"); err != nil {
		return err
	}
	r.Log.Info("Deleting credentials from the store", "name", "rasactl-token", "namespace", r.Namespace)
	return credsStore.Delete("rasactl-token")
}

func (r *RasaCtl) getAuthToken() (string, error) {
	var token string

	credsStore := credentials.Credentials{
		Namespace: r.Namespace,
		Helper:    helpers.Helper,
	}

	r.Log.V(1).Info("Getting credentials from the store", "name", "rasactl-token", "namespace", r.Namespace)
	_, token, err := credsStore.Get("rasactl-token")
	if err != nil {
		return token, fmt.Errorf("%w, use the 'rasa auth login' command", err)
	}

	if !r.RasaXClient.ValidateToken(token) {
		r.initRasaXClient()
		r.Log.V(1).Info("Getting credentials from the store", "name", "rasactl-login", "namespace", r.Namespace)
		username, password, err := credsStore.Get("rasactl-login")
		if err != nil {
			return "", fmt.Errorf("%w, use the 'rasa auth login' command", err)
		}

		authRes, err := r.RasaXClient.Auth(username, password)
		if err != nil {
			return "", err
		}
		r.Log.V(1).Info("Storing credentials in the store", "name", "rasactl-token", "namespace", r.Namespace)
		if err := credsStore.Set("rasactl-token", r.Namespace, authRes.AccessToken); err != nil {
			return "", err
		}
		return authRes.AccessToken, nil
	}

	return token, nil
}

func (r *RasaCtl) isLogged() bool {
	credsStore := credentials.Credentials{
		Namespace: r.Namespace,
		Helper:    helpers.Helper,
	}

	r.Log.Info("Storing credentials in the store", "name", "rasactl-login", "namespace", r.Namespace)
	user, password, err := credsStore.Get("rasactl-login")
	if err != nil {
		r.Log.V(1).Error(err, "Can't get credentials from the store.")
		return false
	}

	return user != "" && password != ""
}

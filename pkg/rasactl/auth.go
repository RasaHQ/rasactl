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
	"encoding/json"
	"fmt"
	"time"

	"github.com/RasaHQ/rasactl/pkg/credentials"
	"github.com/RasaHQ/rasactl/pkg/credentials/helpers"
	"github.com/RasaHQ/rasactl/pkg/utils"
	"github.com/golang-jwt/jwt"
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

	r.Log.Info("Deleteing credentials from the store", "name", "rasactl-login", "namespace", r.Namespace)
	if err := credsStore.Delete("rasactl-login"); err != nil {
		return err
	}
	r.Log.Info("Deleteing credentials from the store", "name", "rasactl-token", "namespace", r.Namespace)
	if err := credsStore.Delete("rasactl-token"); err != nil {
		return err
	}
	return nil
}

func (r *RasaCtl) isJWTExpired(token string) bool {
	t, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		r.Log.Error(err, "Can't parse a JWT token")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		var tm time.Time
		switch exp := claims["exp"].(type) {
		case float64:
			tm = time.Unix(int64(exp), 0)
		case json.Number:
			v, _ := exp.Int64()
			tm = time.Unix(v, 0)
		}

		now := time.Now()

		return now.Before(tm)
	} else {
		r.Log.Error(err, "Can't parse a JWT token")
	}
	return false
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
		return token, fmt.Errorf("%s, use the 'rasa auth login' command", err)
	}

	if !r.isJWTExpired(token) || !r.RasaXClient.ValidateToken(token) {
		r.initRasaXClient()
		r.Log.V(1).Info("Getting credentials from the store", "name", "rasactl-login", "namespace", r.Namespace)
		username, password, err := credsStore.Get("rasactl-login")
		if err != nil {
			return token, fmt.Errorf("%s, use the 'rasa auth login' command", err)
		}

		authRes, err := r.RasaXClient.Auth(username, password)
		if err != nil {
			return token, err
		}
		r.Log.V(1).Info("Storing credentials in the store", "name", "rasactl-token", "namespace", r.Namespace)
		if err := credsStore.Set("rasactl-token", r.Namespace, authRes.AccessToken); err != nil {
			return token, err
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

	if user != "" && password != "" {
		return true
	}
	return false
}

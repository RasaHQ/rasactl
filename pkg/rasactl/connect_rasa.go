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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/RasaHQ/rasactl/pkg/helm"
	"github.com/RasaHQ/rasactl/pkg/types"
	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasa"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

// ConnectRasa connects a local rasa server to a given deployment.
func (r *RasaCtl) ConnectRasa() error {
	r.Spinner.Message("Connecting Rasa Server to Rasa X")
	rasaToken := uuid.New().String()
	environmentName := "production-worker"

	stateData, err := r.KubernetesClient.ReadSecretWithState()
	if err != nil {
		return err
	}

	configDir := string(stateData[types.StateProjectPath])
	if configDir == "" {
		configDir = fmt.Sprintf("/tmp/rasactl-%s", r.Namespace)

		r.Log.V(1).Info("Creating directory", "dir", configDir)

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			err := os.Mkdir(configDir, 0755)
			if err != nil {
				return err
			}
		}
	}

	fileCreds := fmt.Sprintf("%s/.credentials.yaml", configDir)
	fileEndpoints := fmt.Sprintf("%s/.endpoints.yaml", configDir)

	mutualArgs := []string{
		"run",
		"--verbose",
		"--enable-api",
		"--cors",
		"'*'",
		"--auth-token",
		rasaToken,
		"--credentials",
		fileCreds,
		"--endpoints",
		fileEndpoints,
	}
	var productionPort int = r.Flags.ConnectRasa.Port
	var workerPort int = r.Flags.ConnectRasa.Port

	if r.Flags.ConnectRasa.RunSeparateWorker {
		workerPort = workerPort + 1
		environmentName = "production"
	}

	if len(r.Flags.ConnectRasa.ExtraArgs) != 0 {
		mutualArgs = append(mutualArgs, r.Flags.ConnectRasa.ExtraArgs...)
	}

	r.Log.Info("Connecting Rasa Server to Rasa X")
	if r.KubernetesClient.GetBackendType() == types.KubernetesBackendLocal {
		r.HelmClient.SetValues(utils.MergeMaps(helm.ValuesRabbitMQNodePort(),
			helm.ValuesPostgreSQLNodePort(), helm.ValuesHostNetworkRasaX(), r.HelmClient.GetValues()))
		r.Log.V(1).Info("Merging values", "result", r.HelmClient.GetValues())
	} else {
		return errors.Errorf(
			"It looks like you're not using kind as a backend for Kubernetes cluster, the connect rasa command is available only if you use kind.",
		)
	}

	r.Log.Info("Upgrading configuration for Rasa X deployment")
	if err := r.HelmClient.Upgrade(); err != nil {
		return err
	}

	r.Log.Info("Updating configuration for Rasa X")
	if err := r.KubernetesClient.UpdateRasaXConfig(rasaToken); err != nil {
		return err
	}

	r.Log.Info("Restarting Rasa X pod")
	if err := r.KubernetesClient.DeleteRasaXPods(); err != nil {
		return err
	}

	if err := r.saveRasaCredentialsFile(fileCreds); err != nil {
		return err
	}

	if err := r.saveRasaEndpointsFile(fileEndpoints); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*360)
	defer cancel()
	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	r.initRasaXClient()
	r.RasaXClient.URL = url
	if err := r.RasaXClient.WaitForDatabaseMigration(ctx); err != nil {
		return err
	}

	msg := "Starting Rasa Server"
	r.Spinner.Message(msg)
	r.Log.Info(msg, "args", mutualArgs)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ready := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		ready <- true
	}()

	go func(ready chan bool) {
		done := make(chan bool, 1)
		args := mutualArgs
		args = append(args,
			"-p",
			fmt.Sprintf("%d", productionPort),
		)
		r.runRasaServer(environmentName, args, done, ready)
		<-done
	}(ready)

	if r.Flags.ConnectRasa.RunSeparateWorker {
		r.Log.Info("Running separate Rasa X server for the worker environment")
		go func(ready chan bool) {
			done := make(chan bool, 1)
			args := mutualArgs
			args = append(args,
				"-p",
				fmt.Sprintf("%d", workerPort),
			)
			r.runRasaServer("worker", args, done, ready)
			<-done
		}(ready)
	}
	r.Spinner.Stop()
	<-ready
	fmt.Println("exiting")

	return nil
}

func (r *RasaCtl) runRasaServer(environment string, args []string, done chan bool, ready chan bool) {

	cmd := exec.Command("rasa", args...)
	stdout, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err = cmd.Start(); err != nil {
		panic(err)
	}

	// get real-time output
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Printf("(%s) %s", environment, string(tmp))
		if err != nil {
			break
		}
	}

	fmt.Println("awaiting signal")
	<-ready
	fmt.Println("exiting")
	if err = cmd.Wait(); err != nil {
		panic(err)
	}
	fmt.Println("done", cmd.ProcessState)
	done <- true

}

func (r *RasaCtl) saveRasaCredentialsFile(file string) error {
	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}

	creds := rtypes.CredentialsFile{}
	creds.Rasa.URL = fmt.Sprintf("%s/api", url)

	r.Log.Info("Saving credentials.yaml configuration file", "file", file)

	data, err := yaml.Marshal(&creds)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, data, 0644)
	return err
}

func (r *RasaCtl) saveRasaEndpointsFile(file string) error {
	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}

	token, err := r.GetRasaXToken()
	if err != nil {
		return err
	}

	psqlNodePort, err := r.KubernetesClient.GetPostgreSQLSvcNodePort()
	if err != nil {
		return nil
	}

	rabbitNodePort, err := r.KubernetesClient.GetRabbitMqSvcNodePort()
	if err != nil {
		return nil
	}

	usernamePsql, passwordPsql, err := r.KubernetesClient.GetPostgreSQLCreds()
	if err != nil {
		return nil
	}

	usernameRabbit, passwordRabbit, err := r.KubernetesClient.GetRabbitMqCreds()
	if err != nil {
		return nil
	}

	if err := r.GetAllHelmValues(); err != nil {
		return err
	}

	endpoints := rtypes.EndpointsFile{
		Models: rtypes.EndpointModelSpec{
			URL:                  fmt.Sprintf("%s/api/projects/default/models/tags/production", url),
			Token:                token,
			WaitTimeBetweenPulls: 10,
		},
		TrackerStore: rtypes.EndpointTrackerStoreSpec{
			Type:     "sql",
			Dialect:  "postgresql",
			URL:      "127.0.0.1",
			Port:     psqlNodePort,
			Username: usernamePsql,
			Password: passwordPsql,
			Db:       "tracker",
			LoginDb:  r.HelmClient.GetValues()["global"].(map[string]interface{})["postgresql"].(map[string]interface{})["postgresqlDatabase"].(string),
		},
		EventBroker: rtypes.EndpointEventBrokerSpec{
			Type:     "pika",
			URL:      "127.0.0.1",
			Port:     rabbitNodePort,
			Username: usernameRabbit,
			Password: passwordRabbit,
			Queues:   []string{r.HelmClient.GetValues()["rasa"].(map[string]interface{})["rabbitQueue"].(string)},
		},
	}

	r.Log.Info("Saving endpoints.yaml configuration file", "file", file)

	data, err := yaml.Marshal(&endpoints)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, data, 0644)
	return err
}

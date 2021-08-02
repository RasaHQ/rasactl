package rasaxctl

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/helm"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *RasaXCTL) ConnectRasa() error {
	r.Spinner.Message("Connecting Rasa Server to Rasa X")
	rasaToken := uuid.New().String()
	environmentName := "production-worker"
	mutualArgs := []string{
		"run",
		"--verbose",
		"--enable-api",
		"--cors",
		"'*'",
		"--auth-token",
		rasaToken,
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
	if r.KubernetesClient.BackendType == types.KubernetesBackendLocal {
		r.HelmClient.Values = utils.MergeMaps(helm.ValuesRabbitMQNodePort(), helm.ValuesPostgreSQLNodePort(), helm.ValuesHostNetworkRasaX(), r.HelmClient.Values)
		r.Log.V(1).Info("Merging values", "result", r.HelmClient.Values)
	} else {
		return errors.Errorf("It looks like you're not using kind as a backend for Kubernetes cluster, the connect rasa command is available only if you use kind.")
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
	r.Log.Info(msg)
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

func (r *RasaXCTL) runRasaServer(environment string, args []string, done chan bool, ready chan bool) *exec.Cmd {

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

	return cmd
}

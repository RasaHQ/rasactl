package rasax

import (
	"context"
	"fmt"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/pkg/errors"
)

func (r *RasaX) WaitForDatabaseMigration() error {

	for {
		healthStatus, err := r.GetHealthEndpoint()
		if err != nil {
			return err
		}

		if healthStatus == nil {
			msg := "Waiting for the health endpoint to be reachable"
			r.Log.Info(msg, "health", healthStatus)
			r.SpinnerMessage.Message(msg)
			time.Sleep(time.Second * 5)
			continue
		}

		datadabaseStatus := healthStatus.DatabaseMigration

		if datadabaseStatus.Status != "completed" {
			msg := "Waiting for database migration to be completed"
			r.Log.Info(msg, "health", healthStatus)
			r.SpinnerMessage.Message(fmt.Sprintf("%s...%f", msg, datadabaseStatus.ProgressInPercent))
		} else {
			msg := "Database migration is completed"
			r.Log.Info(msg)
			r.SpinnerMessage.Message(msg)
			return nil
		}
		time.Sleep(time.Second * 5)
	}
}

func (r *RasaX) WaitForRasaXWorker() error {

	for {
		healthStatus, err := r.GetHealthEndpoint()
		if err != nil {
			return err
		}

		if healthStatus == nil {
			msg := "Waiting for the health endpoint to be reachable"
			r.Log.Info(msg, "health", healthStatus)
			r.SpinnerMessage.Message(msg)
			time.Sleep(time.Second * 5)
			continue
		}

		workerStatus := healthStatus.Worker

		if workerStatus.Status != 200 {
			msg := "Waiting for the Rasa worker to be ready"
			r.Log.Info(msg, "health", healthStatus)
			r.SpinnerMessage.Message(fmt.Sprintf("%s, status: %d", msg, workerStatus.Status))
		} else {
			msg := "The Rasa worker is ready"
			r.Log.Info(msg)
			r.SpinnerMessage.Message(msg)
			return nil
		}
		time.Sleep(time.Second * 5)
	}
}

func (r *RasaX) WaitForRasaX() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)
	ready := make(chan bool)
	var returnErr error

	go func() {
		for {
			select {
			case <-ticker.C:
				err := r.WaitForDatabaseMigration()
				networkError, _ := utils.CheckNetworkError(err)
				if err != nil && networkError != utils.NetworkErrorConnectionRefused {
					returnErr = err
				}

				err = r.WaitForRasaXWorker()
				networkError, _ = utils.CheckNetworkError(err)
				if err != nil && networkError != utils.NetworkErrorConnectionRefused {
					returnErr = err
				} else if networkError != utils.NetworkErrorConnectionRefused {
					ready <- true
				} else {
					msg := "Waiting for the Rasa X health endpoint to be reachable"
					r.Log.Info(msg)
					r.SpinnerMessage.Message(msg)
				}
			case <-ctx.Done():
				returnErr = errors.Errorf("Error while waiting for Rasa X, error: %s", ctx.Err())
			}
		}
	}()
	<-ready

	r.Log.Info("Rasa X is healthy")
	return returnErr
}

package rasax

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/pkg/errors"
)

func (r *RasaX) WaitForDatabaseMigration() error {
	healthStatus, err := r.GetHealthEndpoint()
	if err != nil {
		return err
	}

	if healthStatus == nil {
		msg := "Waiting for the health endpoint to be reachable"
		r.Log.Info(msg, "health", healthStatus)
		r.SpinnerMessage.Message(msg)
		return nil
	}

	datadabaseStatus := healthStatus.DatabaseMigration

	if datadabaseStatus.Status != "completed" {
		msg := "Waiting for database migration to be completed"
		r.Log.Info(msg, "health", healthStatus)
		r.SpinnerMessage.Message(fmt.Sprintf("%s...%.2f", msg, datadabaseStatus.ProgressInPercent))
	} else {
		msg := "Database migration is completed"
		r.Log.Info(msg)
		r.SpinnerMessage.Message(msg)
		return nil
	}

	return nil
}

func (r *RasaX) WaitForRasaXWorker() error {

	healthStatus, err := r.GetHealthEndpoint()
	if err != nil {
		return err
	}

	if healthStatus == nil {
		msg := "Waiting for the health endpoint to be reachable"
		r.Log.Info(msg, "health", healthStatus)
		r.SpinnerMessage.Message(msg)
		return nil
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
	return nil
}

func (r *RasaX) WaitForRasaX() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*360)
	defer cancel()

	var returnErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			default:

				returnErr = nil

				errWaitForDatabaseMigration := r.WaitForDatabaseMigration()
				networkErrorerrWaitForDatabaseMigration, _ := utils.CheckNetworkError(errWaitForDatabaseMigration)
				if errWaitForDatabaseMigration != nil && networkErrorerrWaitForDatabaseMigration != utils.NetworkErrorConnectionRefused {
					returnErr = errWaitForDatabaseMigration
				}

				errWaitForRasaXWorker := r.WaitForRasaXWorker()
				networkErrorWaitForRasaXWorker, _ := utils.CheckNetworkError(errWaitForRasaXWorker)
				if errWaitForRasaXWorker != nil && networkErrorWaitForRasaXWorker != utils.NetworkErrorConnectionRefused {
					returnErr = errWaitForRasaXWorker
				}

				if networkErrorWaitForRasaXWorker != utils.NetworkErrorConnectionRefused && networkErrorerrWaitForDatabaseMigration != utils.NetworkErrorConnectionRefused && returnErr == nil {
					return
				} else {
					msg := "Waiting for the Rasa X health endpoint to be reachable"
					r.Log.Info(msg)
					r.SpinnerMessage.Message(msg)
				}
				time.Sleep(time.Second * 5)
			case <-ctx.Done():
				returnErr = errors.Errorf("Error while waiting for Rasa X, error: %s", ctx.Err())
				return
			}
		}
	}(ctx)

	wg.Wait()

	r.Log.Info("Rasa X is healthy")
	return returnErr
}

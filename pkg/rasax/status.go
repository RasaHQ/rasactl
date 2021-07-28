package rasax

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (r *RasaX) WaitForDatabaseMigration(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("Error while waiting for Rasa X Database migration status, error: %s", ctx.Err())
		default:
			healthStatus, err := r.GetHealthEndpoint()
			if err != nil {
				return err
			}

			if healthStatus == nil {
				msg := "Waiting for the health endpoint to be reachable"
				r.Log.Info(msg, "health", healthStatus)
				r.SpinnerMessage.Message(msg)
			}

			datadabaseStatus := healthStatus.DatabaseMigration
			if datadabaseStatus.Status != "completed" {
				msg := "Waiting for database migration to be completed"
				r.Log.Info(msg, "health", healthStatus)
				r.SpinnerMessage.Message(fmt.Sprintf("%s...%.2f", msg, datadabaseStatus.ProgressInPercent))
			} else if healthStatus != nil {
				msg := "Database migration is completed"
				r.Log.Info(msg)
				r.SpinnerMessage.Message(msg)
				return nil
			}
			time.Sleep(time.Second * 5)
		}
	}
}

func (r *RasaX) WaitForRasaXWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("Error while waiting for Rasa worker status, error: %s", ctx.Err())
		default:

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
			} else if healthStatus != nil {
				msg := "The Rasa worker is ready"
				r.Log.Info(msg)
				r.SpinnerMessage.Message(msg)
				return nil
			}
			time.Sleep(time.Second * 5)
		}
	}
}

func (r *RasaX) WaitForRasaX() error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*360)
	defer cancel()
	eg, ctx := errgroup.WithContext(c)

	eg.Go(func() error {
		return r.WaitForDatabaseMigration(ctx)
	})

	eg.Go(func() error {
		return r.WaitForRasaXWorker(ctx)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	r.Log.Info("Rasa X is healthy")
	return nil
}

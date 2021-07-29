package rasaxctl

import "github.com/RasaHQ/rasaxctl/pkg/utils"

func (r *RasaXCTL) Upgrade() error {

	if err := utils.ValidateName(r.HelmClient.Namespace); err != nil {
		return err
	}

	// Init Rasa X client
	r.initRasaXClient()

	r.Spinner.Message("Upgrading Rasa X")
	if err := r.HelmClient.Upgrade(); err != nil {
		return err
	}

	if err := r.GetAllHelmValues(); err != nil {
		return err
	}

	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	r.RasaXClient.URL = url

	if err := r.RasaXClient.WaitForRasaX(); err != nil {
		return err
	}

	rasaXVersion, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	helmRelease, err := r.HelmClient.GetStatus()
	if err != nil {
		return err
	}

	if err := r.KubernetesClient.UpdateSecretWithState(rasaXVersion, helmRelease); err != nil {
		return err
	}

	return nil
}

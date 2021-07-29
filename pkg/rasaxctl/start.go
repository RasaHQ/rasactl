package rasaxctl

import "github.com/RasaHQ/rasaxctl/pkg/utils"

func (r *RasaXCTL) Start() error {

	if err := utils.ValidateName(r.HelmClient.Namespace); err != nil {
		return err
	}

	if err := r.KubernetesClient.CreateNamespace(); err != nil {
		return err
	}

	if err := r.KubernetesClient.AddNamespaceLabel(); err != nil {
		return err
	}

	// Init Rasa X client
	r.initRasaXClient()

	if err := r.startOrInstall(); err != nil {
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

	token, err := r.GetRasaXToken()
	if err != nil {
		return err
	}
	r.RasaXClient.Token = token

	if err := r.checkDeploymentStatus(); err != nil {
		return err
	}

	return nil
}

package rasaxctl

import "fmt"

func (r *RasaXCTL) Add() error {
	r.Log.Info("Adding existing project", "namespace", r.Namespace, "releaseName", r.HelmClient.Configuration.ReleaseName)

	release, err := r.HelmClient.GetStatus()
	if err != nil {
		return err
	}

	r.initRasaXClient()
	url, err := r.GetRasaXURL()
	if err != nil {
		return err
	}
	r.RasaXClient.URL = url

	rasaXVersion, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}
	if err := r.KubernetesClient.SaveSecretWithState(); err != nil {
		return err
	}
	if err := r.KubernetesClient.UpdateSecretWithState(rasaXVersion, release); err != nil {
		return err
	}

	if err := r.KubernetesClient.AddNamespaceLabel(); err != nil {
		return err
	}

	fmt.Printf("The %s namespace has been added as a project\n", r.Namespace)

	return nil
}

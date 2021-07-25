package rasaxctl

import (
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/helm"
	"github.com/RasaHQ/rasaxctl/pkg/k8s"
	"github.com/RasaHQ/rasaxctl/pkg/logger"
	"github.com/RasaHQ/rasaxctl/pkg/rasax"
	"github.com/RasaHQ/rasaxctl/pkg/status"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type RasaXCTL struct {
	KubernetesClient *k8s.Kubernetes
	HelmClient       *helm.Helm
	RasaXClient      *rasax.RasaX
	Log              logr.Logger
	Spinner          *status.SpinnerMessage
	Namespace        string
	isRasaXRunning   bool
	isRasaXDeployed  bool
}

func (r *RasaXCTL) InitClients() error {
	r.Log = logger.New()
	r.Spinner = &status.SpinnerMessage{}
	r.Spinner.New()

	r.KubernetesClient = &k8s.Kubernetes{
		Namespace: r.Namespace,
		Log:       r.Log,
	}
	if err := r.KubernetesClient.New(); err != nil {
		return err
	}

	r.HelmClient = &helm.Helm{
		Log:       r.Log,
		Namespace: r.Namespace,
		Spinner:   r.Spinner,
	}
	if err := r.HelmClient.New(); err != nil {
		return err
	}
	r.HelmClient.KubernetesBackendType = r.KubernetesClient.BackendType

	return nil
}

func (r *RasaXCTL) CheckDeploymentStatus() (bool, bool, error) {
	// Check if a Rasa X deployment is already installed and running
	isRasaXDeployed, err := r.HelmClient.IsDeployed()
	if err != nil {
		return false, false, err
	}
	r.isRasaXDeployed = isRasaXDeployed

	isRasaXRunning, err := r.KubernetesClient.IsRasaXRunning()
	if err != nil {
		return false, false, err
	}

	r.isRasaXRunning = isRasaXRunning

	return isRasaXDeployed, isRasaXRunning, nil
}

func (r *RasaXCTL) Start() error {

	if err := utils.ValidateName(r.HelmClient.Namespace); err != nil {
		return err
	}

	if err := r.KubernetesClient.CreateNamespace(); err != nil {
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

func (r *RasaXCTL) Stop() error {
	r.Spinner.Message("Stopping Rasa X")
	if err := r.KubernetesClient.ScaleDown(); err != nil {
		return err
	}
	r.Spinner.Message(fmt.Sprintf("Rasa X for the %s project has been stopped", r.Namespace))
	r.Spinner.Stop()
	return nil
}

func (r *RasaXCTL) startOrInstall() error {
	// Install Rasa X
	if !r.isRasaXDeployed && !r.isRasaXRunning {
		r.Spinner.Message("Deploying Rasa X")
		if viper.GetString("project-path") != "" {
			volume, err := r.KubernetesClient.CreateVolume(viper.GetString("project-path"))
			if err != nil {
				return err
			}
			r.HelmClient.PVCName = volume
		}

		err := r.HelmClient.Install()
		if err != nil {
			return err
		}
	} else if !r.isRasaXRunning {
		// Start Rasa X if deployments are scaled down to 0
		msg := "Starting Rasa X"
		r.Spinner.Message(msg)
		r.Log.Info(msg)

		// Set configuration used for starting a stopped project.
		r.HelmClient.Configuration.StartProject = true

		err := r.HelmClient.Upgrade()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RasaXCTL) GetAllHelmValues() error {
	allValues, err := r.HelmClient.GetValues()
	if err != nil {
		return err
	}
	r.KubernetesClient.Helm.Values = allValues

	return nil
}

func (r *RasaXCTL) GetRasaXURL() (string, error) {
	url, err := r.KubernetesClient.GetRasaXURL()
	if err != nil {
		return url, err
	}
	r.Log.V(1).Info("Get Rasa X URL", "url", url)
	return url, nil
}

func (r *RasaXCTL) GetRasaXToken() (string, error) {
	token, err := r.KubernetesClient.GetRasaXToken()
	if err != nil {
		return token, err
	}

	return token, nil
}

func (r *RasaXCTL) initRasaXClient() {
	r.RasaXClient = &rasax.RasaX{
		Log:            r.Log,
		SpinnerMessage: r.Spinner,
		WaitTimeout:    r.HelmClient.Configuration.Timeout,
	}
	r.RasaXClient.New()
}

func (r *RasaXCTL) checkDeploymentStatus() error {
	err := r.RasaXClient.WaitForRasaX()
	if err != nil {
		return err
	}
	r.Spinner.Message("Ready!")

	if utils.IsDebugOrVerboseEnabled() {
		r.Log.Info("Rasa X is ready", "url", r.RasaXClient.URL)
	}
	r.Spinner.Stop()

	rasaXVersion, err := r.RasaXClient.GetVersionEndpoint()
	if err != nil {
		return err
	}

	if !r.isRasaXDeployed && !r.isRasaXRunning {
		// Print the status box only if it's a new Rasa X deployment
		status.PrintRasaXStatus(rasaXVersion, r.RasaXClient.URL)
	}
	return nil
}

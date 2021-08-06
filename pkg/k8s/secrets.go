package k8s

import (
	"context"
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	rtypes "github.com/RasaHQ/rasaxctl/pkg/types/rasax"
	"helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const secretName string = "rasaxctl"

func (k *Kubernetes) SaveSecretWithState(projectPath string) error {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Type: "rasa.com/rasaxctl.state",
		Data: map[string][]byte{
			types.StateSecretProjectPath:     []byte(projectPath),
			types.StateSecretHelmReleaseName: []byte(k.Helm.ReleaseName),
		},
	}

	k.Log.Info("Saving secret with the project state", "secret", secret.Name, "namespace", k.Namespace)

	_, err := k.clientset.CoreV1().Secrets(k.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) UpdateSecretWithState(data ...interface{}) error {
	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for _, d := range data {
		switch t := d.(type) {
		case *rtypes.VersionEndpointResponse:
			secret.Data[types.StateSecretRasaXVersion] = []byte(t.RasaX)
			secret.Data[types.StateSecretRasaWorkerVersion] = []byte(t.Rasa.Worker)

			enterprise := "inactive"
			if t.Enterprise {
				enterprise = "active"
			}
			secret.Data[types.StateSecretEnterprise] = []byte(enterprise)

		case *release.Release:
			secret.Data[types.StateSecretHelmChartName] = []byte(t.Chart.Name())
			secret.Data[types.StateSecretHelmChartVersion] = []byte(t.Chart.Metadata.Version)
			secret.Data[types.StateSecretHelmReleaseName] = []byte(t.Name)
			secret.Data[types.StateSecretHelmReleaseStatus] = []byte(t.Info.Status)
		}
	}
	k.Log.Info("Updating secret with the project state", "secret", secret.Name, "namespace", k.Namespace, "data", secret.Data)

	if _, err := k.clientset.CoreV1().Secrets(k.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) ReadSecretWithState() (map[string][]byte, error) {

	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return secret.Data, err
	}

	return secret.Data, nil
}

func (k *Kubernetes) DeleteSecretWithState() error {
	if err := k.clientset.CoreV1().Secrets(k.Namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *Kubernetes) GetPostgreSQLCreds() (string, string, error) {
	secretName := fmt.Sprintf("%s-postgresql", k.Helm.ReleaseName)
	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	username := k.Helm.Values["global"].(map[string]interface{})["postgresql"].(map[string]interface{})["postgresqlUsername"].(string)

	return username, string(secret.Data["postgresql-password"]), nil
}

func (k *Kubernetes) GetRabbitMqCreds() (string, string, error) {
	secretName := fmt.Sprintf("%s-rabbit", k.Helm.ReleaseName)
	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	username := k.Helm.Values["rabbitmq"].(map[string]interface{})["rabbitmq"].(map[string]interface{})["username"].(string)

	return username, string(secret.Data["rabbitmq-password"]), nil
}

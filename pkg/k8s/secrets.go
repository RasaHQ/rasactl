package k8s

import (
	"context"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const secretName string = "rasaxctl"

func (k *Kubernetes) SaveSecretWithState() error {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Type: "rasa.com/rasaxctl.state",
		Data: map[string][]byte{
			"project-path": []byte(viper.GetString("project-path")),
		},
	}

	k.Log.Info("Saving secret with the project state", "secret", secret.Name, "namespace", k.Namespace)

	_, err := k.clientset.CoreV1().Secrets(k.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) UpdateSecretWithState(data interface{}) error {
	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	switch t := data.(type) {
	case *types.VersionEndpointResponse:
		secret.Data["rasa-x-version"] = []byte(t.RasaX)
		secret.Data["rasa-worker-version"] = []byte(t.Rasa.Worker)

		enterprise := "inactive"
		if t.Enterprise {
			enterprise = "active"
		}
		secret.Data["enterprise"] = []byte(enterprise)
	}

	k.Log.Info("Updating secret with the project state", "secret", secret.Name, "namespace", k.Namespace, "data", data)

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

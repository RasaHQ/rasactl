package k8s

import (
	"context"
	"fmt"

	types "github.com/RasaHQ/rasaxctl/pkg/types/rasax"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) CreateRasaXConfig() error {

	configSpec := types.EnvironmentsConfigurationFile{
		Rasa: types.RasaSpecEnvironments{
			Production: types.EnvironmentsConfigurationSpec{
				Url:   "",
				Token: "",
			},
		},
	}

	configData, err := yaml.Marshal(&configSpec)
	if err != nil {
		return err
	}

	config := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", k.Helm.ReleaseName, types.RasaXKubernetesConfigMapName),
			Labels: map[string]string{
				"rasaxctl": "true",
			},
		},
		Data: map[string]string{
			"environments": string(configData),
		},
	}

	_, err = k.clientset.CoreV1().ConfigMaps(k.Namespace).Create(context.TODO(), config, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

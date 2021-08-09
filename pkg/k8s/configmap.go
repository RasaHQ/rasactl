/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package k8s

import (
	"context"
	"fmt"

	types "github.com/RasaHQ/rasactl/pkg/types/rasax"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) UpdateRasaXConfig(token string) error {

	var productionPort int = k.Flags.ConnectRasa.Port
	var workerPort int = k.Flags.ConnectRasa.Port

	if k.Flags.ConnectRasa.RunSeparateWorker {
		workerPort = workerPort + 1
	}
	urlProduction := fmt.Sprintf("http://gateway.docker.internal:%d", productionPort)
	urlWorker := fmt.Sprintf("http://gateway.docker.internal:%d", workerPort)

	configSpec := types.EnvironmentsConfigurationFile{
		Rasa: types.RasaSpecEnvironments{
			Production: types.EnvironmentsConfigurationSpec{
				Url:   urlProduction,
				Token: token,
			},
			Worker: types.EnvironmentsConfigurationSpec{
				Url:   urlWorker,
				Token: token,
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
				"rasactl": "true",
			},
		},
		Data: map[string]string{
			"environments": string(configData),
		},
	}

	_, err = k.clientset.CoreV1().ConfigMaps(k.Namespace).Update(context.TODO(), config, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

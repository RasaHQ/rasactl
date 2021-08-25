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
package helm

import "github.com/google/uuid"

func valuesMountHostPath(pvcName string) map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"extraVolumes": []map[string]interface{}{
				{
					"name": "rasa-x-local-path",
					"persistentVolumeClaim": map[string]interface{}{
						"claimName": pvcName,
					},
				},
			},
			"extraVolumeMounts": []map[string]interface{}{
				{
					"name":      "rasa-x-local-path",
					"mountPath": "/project",
				},
			},
		},
	}

	return values
}

func valuesDisableRasaProduction() map[string]interface{} {
	values := map[string]interface{}{
		"rasa": map[string]interface{}{
			"versions": map[string]interface{}{
				"rasaProduction": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}

	return values
}

func valuesDisableRasaWorker() map[string]interface{} {
	values := map[string]interface{}{
		"rasa": map[string]interface{}{
			"versions": map[string]interface{}{
				"rasaWorker": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}

	return values
}

func valuesUseDedicatedKindNode(namespace string) map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"tolerations": []map[string]interface{}{
				{
					"key":      "rasactl",
					"operator": "Equal",
					"value":    "true",
					"effect":   "NoSchedule",
				},
			},
			"nodeSelector": map[string]interface{}{
				"rasactl-project": namespace,
			},
		},
	}

	return values
}

func valuesDisableNginx() map[string]interface{} {

	values := map[string]interface{}{
		"nginx": map[string]interface{}{
			"enabled": false,
		},
	}

	return values
}

func valuesNginxNodePort() map[string]interface{} {

	values := map[string]interface{}{
		"nginx": map[string]interface{}{
			"service": map[string]interface{}{
				"type": "NodePort",
			},
		},
	}

	return values
}

func valuesSetupLocalIngress(host string) map[string]interface{} {
	values := map[string]interface{}{
		"ingress": map[string]interface{}{
			"hosts": []map[string]interface{}{
				{
					"host":  host,
					"paths": []string{"/"},
				},
			},
		},
	}

	return values
}

func valuesSetRasaXPassword(password string) map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"initialUser": map[string]interface{}{
				"password": password,
			},
		},
	}

	return values
}

// ValuesHostNetworkRasaX returns helm values that set hostNetwork to 'true'
// for Rasa X deployment.
func ValuesHostNetworkRasaX() map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"hostNetwork": true,
		},
	}

	return values
}

// ValuesRabbitMQNodePort returns helm values which set the rabbitmq service type to NodePort.
func ValuesRabbitMQNodePort() map[string]interface{} {
	values := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"service": map[string]interface{}{
				"type": "NodePort",
			},
		},
	}

	return values
}

// ValuesPostgreSQLNodePort return helm vales which set the postgresql service type to NodePort.
func ValuesPostgreSQLNodePort() map[string]interface{} {
	values := map[string]interface{}{
		"postgresql": map[string]interface{}{
			"service": map[string]interface{}{
				"type": "NodePort",
			},
		},
	}

	return values
}

func valuesRabbitMQErlangCookie() map[string]interface{} {
	values := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"rabbitmq": map[string]interface{}{
				"erlangCookie": uuid.New().String(),
			},
		},
	}

	return values
}

func valuesUseEdgeReleaseRasaX() map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"tag": "latest",
		},
		"eventService": map[string]interface{}{
			"tag": "latest",
		},
		"dbMigrationService": map[string]interface{}{
			"tag": "latest",
		},
	}

	return values
}

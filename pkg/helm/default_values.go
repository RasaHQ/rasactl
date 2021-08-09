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

func valuesEnableRasaProduction() map[string]interface{} {
	values := map[string]interface{}{
		"rasa": map[string]interface{}{
			"versions": map[string]interface{}{
				"rasaProduction": map[string]interface{}{
					"enabled": true,
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
					"key":      "rasaxctl",
					"operator": "Equal",
					"value":    "true",
					"effect":   "NoSchedule",
				},
			},
			"nodeSelector": map[string]interface{}{
				"rasaxctl-project": namespace,
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

func ValuesHostNetworkRasaX() map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"hostNetwork": true,
		},
	}

	return values
}

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

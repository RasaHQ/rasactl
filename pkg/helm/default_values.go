package helm

func valuesMountHostPath(pvcName string) map[string]interface{} {
	values := map[string]interface{}{
		"rasax": map[string]interface{}{
			"persistence": map[string]interface{}{
				"existingClaim": pvcName,
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

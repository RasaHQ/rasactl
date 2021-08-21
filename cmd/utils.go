package cmd

import (
	"os"
	"syscall"

	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// HandleSignals receives a signal from the channel and runs an action depends on the type of the signal.
func HandleSignals(sigs chan os.Signal) {
	signal := <-sigs
	runOnClose(signal)
}

func runOnClose(signal os.Signal) {
	emoji.Println("Bye :wave:")

	switch signal {
	case os.Interrupt:
		os.Exit(130)
	case syscall.SIGTERM:
		os.Exit(143)
	default:
		os.Exit(0)
	}
}

func parseModelUpDownArgs(namespace, detectedNamespace string, args []string) (string, string, string, error) {
	var modelName, modelPath string
	for {
		switch {
		case namespace == "":
			return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
		case len(args) == 1:
			if args[0] == detectedNamespace {
				return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a model name"))
			} else if detectedNamespace != "" {
				modelName = args[0]
				return modelName, modelPath, detectedNamespace, nil
			} else if detectedNamespace == "" && namespace != "" {
				modelName = args[0]
				return modelName, modelPath, namespace, nil
			} else if detectedNamespace == "" {
				return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a model name"))
			}
		case len(args) == 2 && detectedNamespace != "":
			modelName = args[0]
			modelPath = args[1]
			return modelName, modelPath, detectedNamespace, nil
		case len(args) == 2 && detectedNamespace == "" && args[0] != namespace:
			modelName = args[0]
			modelPath = args[1]
			return modelName, modelPath, namespace, nil
		case len(args) == 2 && detectedNamespace == "" && args[0] == namespace:
			modelName = args[1]
			return modelName, modelPath, namespace, nil
		case len(args) == 3:
			modelName = args[1]
			modelPath = args[2]
			return modelName, modelPath, namespace, nil
		default:
			return "", "", "", nil
		}
	}
}

func parseModelTagArgs(namespace, detectedNamespace string, args []string) (string, string, string, error) {
	var modelName, modelTag string
	for {
		switch {
		case namespace == "":
			return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
		case len(args) == 2:
			if args[0] == detectedNamespace {
				return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a model name"))
			} else if detectedNamespace == "" && namespace != "" {
				modelName = args[0]
				modelTag = args[1]
				return modelName, modelTag, namespace, nil
			} else if detectedNamespace != "" {
				modelName = args[0]
				modelTag = args[1]
				return modelName, modelTag, detectedNamespace, nil
			} else if detectedNamespace == "" {
				return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a tag name"))
			}
		case len(args) == 2 && detectedNamespace != "":
			modelName = args[0]
			modelTag = args[1]
			return modelName, modelTag, detectedNamespace, nil
		case len(args) == 2 && detectedNamespace == "":
			return "", "", "", errors.Errorf(errorPrint.Sprint("Not enough arguments"))
		case len(args) == 3:
			modelName = args[1]
			modelTag = args[2]
			return modelName, modelTag, namespace, nil
		default:
			return "", "", "", nil
		}
	}
}

func setDeploymentIfOnlyOne(cmd *cobra.Command) error {
	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return errors.Errorf(errorPrint.Sprint(err))
	}

	// If there is only one deployment then use it as the default
	if (namespace == "" && len(namespaces) == 1) ||
		(cmd.CalledAs() == "start" && len(namespaces) == 1 && !rasactlFlags.Start.Create) {

		namespace = namespaces[0]
		rasaCtl.Namespace = namespace
		log.Info("Setting default namespace", "namespace", namespace)
		if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
			return err
		}
	}
	return nil
}

func checkIfNamespaceExists() error {
	if namespace == "" {
		return errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
	}

	isNamespaceExist, err := rasaCtl.KubernetesClient.IsNamespaceExist(rasaCtl.Namespace)
	if err != nil {
		return errors.Errorf(errorPrint.Sprintf("%s", err))
	}

	if !isNamespaceExist {
		return errors.Errorf(errorPrint.Sprintf("The %s deployment doesn't exist.\n", rasaCtl.Namespace))
	}
	return nil
}

func parseNamespaceModelDeleteCommand(args []string) error {
	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return errors.Errorf(errorPrint.Sprint(err))
	}
	for {
		switch {
		case len(namespaces) == 1 && len(args) == 1 && args[0] != namespaces[0]:
			namespace = namespaces[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		case len(namespaces) == 1 && args[0] != namespaces[0] && len(args) == 2:
			namespace = args[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		default:
			return nil
		}
	}
}

func parseNamespaceModelDownloadCommand(args []string) error {
	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return errors.Errorf(errorPrint.Sprint(err))
	}
	for {
		switch {
		case len(namespaces) == 1 && len(args) == 1 && args[0] != namespaces[0]:
			namespace = namespaces[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		case len(namespaces) == 1 && args[0] != namespaces[0] && len(args) == 2:
			namespace = namespaces[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		case len(namespaces) == 1 && args[0] != namespaces[0] && len(args) == 3:
			namespace = args[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		default:
			return nil
		}
	}
}

func parseNamespaceModelUploadCommand(args []string) error {
	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return errors.Errorf(errorPrint.Sprint(err))
	}

	for {
		switch {
		case len(namespaces) == 1 && len(args) == 1 && args[0] != namespaces[0]:
			namespace = namespaces[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		case len(namespaces) == 1 && args[0] != namespaces[0] && len(args) == 2:
			namespace = args[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		default:
			return nil
		}
	}
}

func parseNamespaceModelTagCommand(args []string) error {
	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return errors.Errorf(errorPrint.Sprint(err))
	}
	for {
		switch {
		case len(namespaces) == 1 && len(args) == 2 && args[0] != namespaces[0]:
			namespace = namespaces[0]
			rasaCtl.Namespace = namespaces[0]
			if err := rasaCtl.SetNamespaceClients(namespaces[0]); err != nil {
				return err
			}
			return nil
		case len(namespaces) == 1 && args[0] != namespaces[0] && len(args) == 3:
			namespace = args[0]
			rasaCtl.Namespace = namespace
			if err := rasaCtl.SetNamespaceClients(namespace); err != nil {
				return err
			}
			return nil
		default:
			return nil
		}
	}
}

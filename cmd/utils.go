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

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf(
			"%q accepts no arguments\n\nUsage:  %s",
			cmd.CommandPath(),
			cmd.UseLine(),
		)
	}
	return nil
}

func maximumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > n {
			return errors.Errorf(
				"%q accepts at most %d %s\n\nUsage:  %s",
				cmd.CommandPath(),
				n,
				"arguments",
				cmd.UseLine(),
			)
		}
		return nil
	}
}

func parseModelDownloadArgs(namespace, detectedNamespace string, args []string) (string, string, string, error) {
	var modelName, modelPath string
	if namespace == "" {
		return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
	} else if len(args) == 1 {
		if args[0] == detectedNamespace {
			return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a model name"))
		} else if detectedNamespace != "" {
			modelName = args[0]
			return modelName, modelPath, detectedNamespace, nil
		} else if detectedNamespace == "" {
			return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a model name"))
		}
	} else if len(args) == 2 && detectedNamespace != "" {
		modelName = args[0]
		modelPath = args[1]
		return modelName, modelPath, detectedNamespace, nil
	} else if len(args) == 2 && detectedNamespace == "" {
		modelName = args[1]
		return modelName, modelPath, namespace, nil
	} else if len(args) == 3 {
		modelName = args[1]
		modelPath = args[2]
		return modelName, modelPath, namespace, nil
	}

	return "", "", "", nil
}

func parseModelTagArgs(namespace, detectedNamespace string, args []string) (string, string, string, error) {
	var modelName, modelTag string
	if namespace == "" {
		return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a deployment name"))
	} else if len(args) == 2 {
		if args[0] == detectedNamespace {
			return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a model name"))
		} else if detectedNamespace != "" {
			modelName = args[0]
			return modelName, modelTag, detectedNamespace, nil
		} else if detectedNamespace == "" {
			return "", "", "", errors.Errorf(errorPrint.Sprint("You have to pass a tag name"))
		}
	} else if len(args) == 2 && detectedNamespace != "" {
		modelName = args[0]
		modelTag = args[1]
		return modelName, modelTag, detectedNamespace, nil
	} else if len(args) == 2 && detectedNamespace == "" {
		return "", "", "", errors.Errorf(errorPrint.Sprint("Not enough arguments"))
	} else if len(args) == 3 {
		modelName = args[1]
		modelTag = args[2]
		return modelName, modelTag, namespace, nil
	}

	return "", "", "", nil
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

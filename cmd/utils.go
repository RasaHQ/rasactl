package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/docker/docker/pkg/namesgenerator"
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

func parseArgs(args []string, minArgs, maxArgs int) ([]string, error) {
	var isInRange bool = true
	var isMaxArgs bool = false
	var ns string
	var currentNamespace string

	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return nil, errors.Errorf(errorPrint.Sprint(err))
	}

	numNamespaces := len(namespaces)
	// Check if namespace is defined by .rasactl or the configuration file
	if namespace != "" {
		currentNamespace = namespace
	} else if namespace == "" && numNamespaces == 1 {
		currentNamespace = namespaces[0]
	}

	// Check args range
	if len(args) < minArgs || len(args) > maxArgs {
		fmt.Println("dupa range")
		isInRange = false
	}

	// Check if the number of args is equal to maxArgs
	if len(args) == maxArgs {
		isMaxArgs = true
	}

	// Check if a new deployment is requested
	newDeployment := false
	if rasactlFlags.Start.Create || rasactlFlags.Start.Project ||
		rasactlFlags.Start.ProjectPath != "" {
		newDeployment = true
	}

	switch {
	case numNamespaces == 0 && len(args) == 1:
		ns = args[0]
	case numNamespaces == 0 && len(args) == 0:
		ns = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	case numNamespaces >= 0 && len(args) == 0 && newDeployment:
		ns = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	case numNamespaces == 1 && len(args) == 0 && !newDeployment:
		// use default namespace, example: rasactl command
		ns = currentNamespace
	case numNamespaces == 1 && isInRange && isMaxArgs && minArgs != len(args):
		// use default namespace, example: rasactl command arg1
		// the number of args is equal to maxArgs
		ns = args[0]
		args = args[1:]
	case numNamespaces == 1 && isInRange && isMaxArgs && minArgs == len(args):
		// use default namespace, example: rasactl command arg1
		// the number of args is equal to maxArgs
		ns = args[0]
		args = []string{}
	case numNamespaces == 1 && isInRange && !isMaxArgs:
		// use default namespace, example: rasactl command arg1
		// the number of args is not equal to maxArgs
		ns = currentNamespace
	case numNamespaces >= 2 && len(args) == 0 && currentNamespace != "":
		ns = currentNamespace
	case numNamespaces >= 2 && len(args) == 0 && currentNamespace == "":
		ns = ""
	case numNamespaces >= 2 && !isMaxArgs && currentNamespace == "":
		ns = ""
		return nil, fmt.Errorf("can't find default deployment, you have to pass all arguments")
	case numNamespaces >= 2 && isInRange && !isMaxArgs && currentNamespace != "":
		ns = currentNamespace
	case numNamespaces >= 2 && isInRange && isMaxArgs && minArgs != len(args):
		ns = args[0]
		args = args[1:]
	case numNamespaces >= 2 && isInRange && isMaxArgs && minArgs == len(args):
		ns = args[0]
		args = []string{}
	}
	args = append([]string{ns}, args...)

	// extend index up to maxArgs
	for i := 0; i < maxArgs-len(args); i++ {
		args = append(args, "")
	}

	fmt.Printf("current namespace: %s, namespace: %s, args: %s, len: %d\n", currentNamespace, ns, args, len(args))

	// The valid namespace is returned as the first element in the args array
	namespace = ns
	rasaCtl.Namespace = ns
	log.Info("Setting namespace", "namespace", ns)
	if err := rasaCtl.SetNamespaceClients(ns); err != nil {
		return nil, err
	}

	return args, nil
}

package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"

	"github.com/RasaHQ/rasactl/pkg/types"
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

func checkIfDeploymentsExist() error {
	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return errors.Errorf(errorPrint.Sprint(err))
	}

	if len(namespaces) == 0 {
		return errors.Errorf(errorPrint.Sprint("No deployments, use the 'rasactl start' command to create one."))
	}
	return nil
}

//nolint:golint,gocyclo
func parseArgs(currentNamespace string, args []string, minArgs, maxArgs int, flags *types.RasaCtlFlags) ([]string, error) {
	isInRange := true
	isMaxArgs := false
	nsExists := false
	var ns string

	namespaces, err := rasaCtl.KubernetesClient.GetNamespaces()
	if err != nil {
		return nil, errors.Errorf(errorPrint.Sprint(err))
	}

	numNamespaces := len(namespaces)
	// Check if namespace is defined by .rasactl or the configuration file
	if currentNamespace == "" && numNamespaces == 1 {
		currentNamespace = namespaces[0]
	}

	// Check args range
	if len(args) < minArgs || len(args) > maxArgs {
		isInRange = false
	}

	// Check if the number of args is equal to maxArgs
	if len(args) == maxArgs {
		isMaxArgs = true
	}

	// Check if a new deployment is requested
	newDeployment := false
	if flags.Start.Create || flags.Start.Project ||
		flags.Start.ProjectPath != "" {
		newDeployment = true
	}

	// Check if args[0] is a namespace
	if len(args) != 0 {
		nsExists, err = rasaCtl.KubernetesClient.IsNamespaceExist(args[0])
		if err != nil {
			return nil, err
		}
	}

	switch {
	case numNamespaces == 0 && len(args) == 1:
		ns = args[0]
	case numNamespaces == 0 && len(args) == 0:
		ns = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	case numNamespaces >= 0 && len(args) == 0 && newDeployment:
		ns = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	case numNamespaces == 1 && len(args) == 0 && !newDeployment:
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
	case numNamespaces == 1 && isInRange && !isMaxArgs && currentNamespace != "" && nsExists:
		ns = args[0]
		args = args[1:]
	case numNamespaces == 1 && isInRange && !isMaxArgs && currentNamespace != "" && !nsExists:
		ns = currentNamespace
	case numNamespaces >= 2 && len(args) == 0 && currentNamespace != "":
		ns = currentNamespace
	case numNamespaces >= 2 && isInRange && !isMaxArgs && currentNamespace != "" && nsExists:
		ns = args[0]
		args = args[1:]
	case numNamespaces >= 2 && isInRange && !isMaxArgs && currentNamespace != "" && !nsExists:
		ns = currentNamespace
	case numNamespaces >= 2 && (minArgs+1 == len(args)) && currentNamespace == "":
		ns = args[0]
		args = args[1:]
	case numNamespaces >= 2 && len(args) == 0 && currentNamespace == "":
		ns = ""
	case numNamespaces >= 2 && !isMaxArgs && currentNamespace == "":
		ns = ""
	case numNamespaces >= 2 && isInRange && isMaxArgs && minArgs == len(args):
		ns = args[0]
		args = []string{}
	case numNamespaces >= 2 && isInRange && !isMaxArgs && currentNamespace != "":
		ns = currentNamespace
	case numNamespaces >= 2 && isInRange && isMaxArgs && minArgs != len(args):
		ns = args[0]
		args = args[1:]
	default:
		return nil, fmt.Errorf("can't parse arguments")

	}
	args = append([]string{ns}, args...)

	// extend index up to maxArgs
	for i := 0; i < maxArgs-len(args); i++ {
		args = append(args, "")
	}

	// The valid namespace is returned as the first element in the args array
	namespace = ns
	rasaCtl.Namespace = namespace
	log.Info("Setting namespace", "namespace", ns)
	if err := rasaCtl.SetNamespaceClients(ns); err != nil {
		return nil, err
	}

	return args, nil
}

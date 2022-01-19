# rasactl

`rasactl` deploys Rasa X / Enterprise on your local or remote Kubernetes cluster and manages Rasa X / Enterprise deployments.

<p align="center"><img src="img/render1631526529716.gif?raw=true"/></p>

## Features

- deploy Rasa X / Enterprise

  You can use `rasactl` to deploy Rasa X / Enterprise on your local machine or a VM in one of the major cloud providers.

  (check the [Prerequisites](#Prerequisites) section)

- upgrade Rasa X / Enterprise

  Upgrade/change configuration for an existing Rasa X deployment.

- stop/delete/start a running Rasa X / Enterprise deployment

  Manage the lifecycle of your deployment: you can stop, delete or start one of the Rasa X deployments managed by `rasactl`.

- connect a local Rasa Server to Rasa X / Enterprise

  You can use your local Rasa Open Source server along with Rasa X / Enterprise. `rasactl` will prepare configuration for Rasa OSS and Rasa X and run the Rasa Open Source server on your local machine.

  (requires `kind` and Rasa OSS installed locally)

- use a local Rasa project along Rasa X / Enterprise deployment

  Use your local Rasa project along with Rasa X / Enterprise deployment. The `rasactl` provides an easy way to use your local Rasa project along with Rasa X / Enterprise.

  This setup was previously referred to as "local mode" in older Rasa X versions.

  (requires `kind` and Rasa X >= 1.0.0)

## Table of Contents

- [rasactl](#rasactl)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [Linux / macOS](#linux--macos)
  - [Compatibility matrix](#compatibility-matrix)
  - [Before you start](#before-you-start)
  - [Values File](#values-file)
  - [Configuration](#configuration)
    - [Environment variables](#environment-variables)
    - [Configuration file](#configuration-file)
  - [Global flags](#global-flags)
  - [Commands](#commands)
    - [The `add` command](#the-add-command)
    - [The `start` command](#the-start-command)
    - [The `stop` command](#the-stop-command)
    - [The `delete` command](#the-delete-command)
    - [The `list` command](#the-list-command)
    - [The `status` command](#the-status-command)
    - [The `config use-deployment` command](#the-config-use-deployment-command)
    - [The `connect rasa` command](#the-connect-rasa-command)
    - [The `auth login` command](#the-auth-login-command)
    - [The `auth logout` command](#the-auth-logout-command)
    - [The `logs` command](#the-logs-command)
  - [Enterprise Management Commands](#enterprise-management-commands)
    - [The `enterprise activate` command](#the-enterprise-activate-command)
    - [The `enterprise deactivate` command](#the-enterprise-deactivate-command)
  - [Model Management Commands](#model-management-commands)
    - [The `model delete` command](#the-model-delete-command)
    - [The `model download` command](#the-model-download-command)
    - [The `model list` command](#the-model-list-command)
    - [The `model tag` command](#the-model-tag-command)
    - [The `model upload` command](#the-model-upload-command)
    - [Upload a model to Rasa X](#upload-a-model-to-rasa-x)
  - [Examples of usage](#examples-of-usage)
    - [Run Rasa X / Enterprise with a local Rasa Server](#run-rasa-x--enterprise-with-a-local-rasa-server)
    - [Run Rasa X / Enterprise with mounted a local Rasa project](#run-rasa-x--enterprise-with-mounted-a-local-rasa-project)
    - [Upgrade Rasa X / Enterprise version](#upgrade-rasa-x--enterprise-version)
    - [Deploy Rasa X in one of the public cloud providers](#deploy-rasa-x-in-one-of-the-public-cloud-providers)
  - [Development](#development)
    - [How to run it?](#how-to-run-it)
    - [Run unit tests](#run-unit-tests)
    - [Kind cluster for developing purposes](#kind-cluster-for-developing-purposes)
  - [License](#license)

## Prerequisites

- Kubernetes cluster >= 1.14+

   or

- kind (for local mode)

(You can use the [REI](https://github.com/RasaHQ/REI) to install all required components on your local machine or a VM.)

## Installation

### Linux / macOS

- Binary downloads of `rasactl` can be found on [the Releases page](https://github.com/rasahq/rasactl/releases/latest). You can manually install rasactl by coping the binary into your `bin`:

```text
$ curl -L https://github.com/RasaHQ/rasactl/releases/download/0.0.24/rasactl_0.0.24_darwin_amd64.tar.gz -O
$ tar -zxvf rasactl_0.0.24_darwin_amd64.tar.gz
$ cp rasactl_0.0.24_darwin_amd64/rasactl /usr/local/bin/
```

- You can also install via `brew`:

```text
$ brew tap rasahq/rasactl
$ brew install rasactl
```

## Compatibility matrix

| rasactl version | Helm chart version |
| --------------- | ------------------ |
| `>= 1.0.x`      | `4.x`              |
| `<= 0.5.x`      | `3.x`              |

## Before you start

Below you can find several things that are good to know and keep in mind when you use `rasactl`.

- It is possible to configure multiple deployments with `rasactl`. A `rasactl` command will always execute an operation on a single deployment. Here is the order in which `rasactl` determines which deployment to use:

  1. A deployment name passed as an argument in CLI, e.g. `rasactl status deployment-name`, you can use `rasactl help command` to see usage example for a given command.
  2. `rasactl` checks if a `.rasactl` file exists in a current working directory. If so, the deployment defined in the file is used. This `.rasactl` file is created automatically when the `rasactl start --project` command is executed.
  3. `rasactl` checks if a default deployment is configured in the `rasactl.yaml` configuration file, if yes, then the default deployment is used. The default deployment can be set by using the `rasactl config use-deployment` command.
  4. If there is only one deployment, then it's used.

  You can use the [`rasactl list`](#the-list-command) command to check which deployment is used as the current one.

  The `rasactl delete` command requires explicitly passing a deployment name as an argument.

- `rasactl` uses the [`rasa-x-helm`](https://github.com/RasaHQ/rasa-x-helm) chart to deploy Rasa X / Enterprise.
- `rasactl` deploys Rasa X / Enterprise without a Rasa Open Source server. It's up to you to connect Rasa OSS with Rasa X / Enterprise deployment.
- `rasactl` uses a Kubernetes context from the kubeconfig file, if you want to switch Kubernetes cluster you have to use `kubectl` or other tools that change the active context for the kubeconfig.

## Values File

The `rasactl` uses the [`rasa-x-helm` chart](https://github.com/RasaHQ/rasa-x-helm) to deploy Rasa X / Enterprise, which means you can use [the helm chart values](https://github.com/RasaHQ/rasa-x-helm/blob/main/charts/rasa-x/values.yaml) to configure deployment. The `rasactl` enables template usage for the values file so that it's possible to use the [Go template](https://pkg.go.dev/text/template#hdr-Actions) and [Sprig function](http://masterminds.github.io/sprig/) within the value file, e.g.

```yaml
# values.yaml
rasax:
  podLabels:
    rasactl: "true"
    test_version: {{ env "RASACTL_TEST_VERSION" }}
    test_template: {{ coalesce 0 1 2 }}
```

## Configuration

### Environment variables

|                  Name                  |                                                                                                                         Description                                                                                                                          |
| -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `RASACTL_AUTH_USER`                    | The username that is used to authorize to Rasa X / Enterprise                                                                                                                                                                                                |
| `RASACTL_AUTH_PASSWORD`                | The password that is used to authorize to Rasa X / Enterprise                                                                                                                                                                                                |
| `RASACTL_RASA_X_URL`                   | Set Rasa X / Enterprise URL. By default, the URL is detected automatically, but if you use a custom configuration and you wanna define Rasa X URL explicitly you can use the env variable. The `RASACTL_RASA_X_URL` overrides Rasa X URL for all deployment. |
| `RASACTL_RASA_X_URL_<DEPLOYMENT_NAME>` | Set Rasa X / Enterprise URL for a given deployment, e.g. if a deployment name is `my-deployment`, then you can use the `RASACTL_RASA_X_URL_MY_DEPLOYMENT` environment variable to define the Rasa X URL for the `my-deployment`.                             |
| `RASACTL_KUBECONFIG`                   | Absolute path to the kubeconfig file (default "`$HOME/.kube/config`")                                                                                                                                                                                        |

### Configuration file

Below you can find an example of the configuration file and parameters that can be defined, by default configuration file is located in `$HOME/.rasactl.yaml`.

```yaml
# Deployment name that is used as a current deployment (default).
# You can use the `rasactl config use-deployment` command to set the current deployment.
current-deployment: my-deployment

# Name of the kubeconfig context to use
kube-context: ""

# Absolute path to the kubeconfig file
kubeconfig: /home/user/.kube/config
```

## Global flags

Below you can find global flags that can be used with every command.

```text
Global Flags:
      --config string         config file (default is $HOME/.rasactl.yaml)
      --debug                 enable debug output
  -h, --help                  help for rasactl
      --kube-context string   name of the kubeconfig context to use
      --kubeconfig string     absolute path to the kubeconfig file (default "$HOME/.kube/config")
      --verbose               enable verbose output
```

## Commands

```text
Available Commands:
  add         add existing Rasa X deployment to rasactl
  auth        manage credentials for Rasa X / Enterprise
  completion  generate the autocompletion script for the specified shell
  config      modify the configuration file
  connect     connect a component (e.g. a Rasa OSS server) to Rasa X
  delete      delete Rasa X deployment
  enterprise  manage Rasa Enterprise
  help        Help about any command
  list        list deployments
  logs        print the logs for a container in a pod
  model       manage models for Rasa X / Enterprise
  open        open Rasa X in a web browser
  start       start a Rasa X deployment
  status      show deployment status
  stop        stop Rasa X deployment
  upgrade     upgrade Rasa X deployment
```

### The `add` command

Adds existing Rasa X deployment to rasactl.

If you already have a Rasa X deployment that uses the rasa-x-helm chart you can add the deployment and manage it by rasactl.

```text
Usage:
  rasactl add NAMESPACE [flags]
```

```text
Examples:
  # Add a Rasa X deployment that is deployed in the 'my-test' namespace.
  $ rasactl add my-test

  # Add a Rasa X deployment that is deployed in the 'my-test' namespace and
  # a helm release name for the deployment is 'rasa-x-example'.
  $ rasactl add my-test --rasa-x-release-name rasa-x-example
```

```text
Flags:
  -h, --help                         help for add
      --rasa-x-release-name string   a helm release name to manage (default "rasa-x")
```

### The `start` command

The `start` command creates a Rasa X deployment or starts a stopped deployment if a given deployment already exists.

```text
Usage:
  rasactl start [DEPLOYMENT-NAME] [flags]

```

```text
Examples:
  # Create a new Rasa X deployment with an autogenerated name.
  $ rasactl start

  # Create a Rasa X deployment with a defined name.
  $ rasactl start my-deployment

  # Create a new deployment if there is already one or more deployments.
  # rasactl start --create

  # Create a Rasa X deployment with custom configuration, e.g the following configuration changes a Rasa X version.
  # All available values: https://github.com/RasaHQ/rasa-x-helm/blob/main/charts/rasa-x/values.yaml
  $ rasactl start --values-file custom-configuration.yaml

  # Create a Rasa X deployment with a defined password.
  $ rasactl start --rasa-x-password mypassword

  # Create a Rasa X deployment that uses a local Rasa project.
  # The command is executed in a Rasa project directory.
  $ rasactl start --project
```

```text
Flags:
      --create                        create a new deployment. If --project or --project-path is set, or there is no existing deployment, the flag is not required to create a new deployment
  -h, --help                          help for start
  -p, --project                       use the current working directory as a project directory, the flag is ignored if --project-path is used
      --project-path string           absolute path to the project directory mounted in kind
      --rasa-x-chart-version string   a helm chart version to use
      --rasa-x-edge-release           use the latest edge release of Rasa X
      --rasa-x-password string        Rasa X password (default "rasaxlocal")
      --rasa-x-password-stdin         read the Rasa X password from stdin
      --rasa-x-release-name string    a helm release name to manage (default "rasa-x")
      --values-file string            absolute path to the values file
      --wait-timeout duration         time to wait for Rasa X to be ready (default 10m0s)
```

### The `stop` command

The `stop` command stops a running Rasa X / Enterprise deployment. The Rasa X deployment and all its components will be scaled down to 0.

```text
Usage:
  rasactl stop [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Stop a Rasa X deployment with the 'my-deployment' name.
  $ rasactl stop my-deployment

  # Stop a currently active Rasa X deployment.
  # The command stops the currently active deployment.
  # You can use the 'rasactl list' command to check which deployment is currently used.
  $ rasactl stop
```

```text
Flags:
  -h, --help   help for stop
```

### The `delete` command

The `delete` command deletes a Rasa X / Enterprise deployment.

You can use the `--prune` flag to remove a namespace where Rasa X deployment is located.

***Notice*** If you want to free resources, or temporarily you don't need to run Rasa X deployment, you can stop a Rasa X / Enterprise deployment instead of deleting it. Stopping the Rasa X deployment will free resources, but keep the current configuration.

```text
Usage:
  rasactl delete DEPLOYMENT-NAME [flags]

Aliases:
  delete, del
```

```text
Examples:
  # Delete the 'my-example' deployment.
  $ rasactl delete my-example

  # Prune the 'my-example' deployment, execute the command with the --prune flag deletes the whole namespace.
  $ rasactl delete my-example --prune
```

```text
Flags:
      --force   if true, delete resources and ignore errors
  -h, --help    help for delete
      --prune   if true, delete a namespace with a project
```

### The `list` command

List all deployments.

```text
$ rasactl list
CURRENT	NAME         	STATUS 	RASA PRODUCTION	RASA WORKER	ENTERPRISE	VERSION
       	hopeful-haibt	Running	2.8.1          	2.8.1      	inactive  	0.42.0
*      	vibrant-yalow	Running	2.8.1          	2.8.1      	inactive  	0.42.0
```

The `*` in the `CURRENT` field indicates a deployment that is used as default. It means that every time when you execute `rasactl` command without defining the deployment name, the deployment marked with `*` is used.

A deployment is marked as `CURRENT` if:

- there is a `.rasactl` file that includes a deployment name in your current working directory. The file is automatically created if you run the `rasactl start` command with the `--project` or `--project-path` flag
- there is only one deployment
- you set the current deployment by using the `rasactl config use-deployment` command

### The `status` command

Show the status of a deployment.

```text
Usage:
  rasactl status [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Show status for the 'example' deployment.
  $ rasactl status example

  # Show status for the 'example' deployment along with details.
  $ rasactl status example --details
```

```text
Flags:
  -d, --details         show detailed information, such as running pods, helm chart status
  -h, --help            help for status
  -o, --output string   output format. One of: json|table (default "table")
```

Example output:

```text
$ rasactl status vibrant-yalow
Name:                   	vibrant-yalow
Status:                 	Running
URL:                    	http://vibrant-yalow.rasactl.localhost
Version:                	0.42.0
Enterprise:             	inactive
Rasa production version:	2.8.1
Rasa worker version:    	2.8.1
Project path:           	/home/ubuntu/test
```

### The `config use-deployment` command

Sets the current-deployment in the configuration file.

If you have multiple Rasa X deployments, and you are not in a project directory you have to explicitly define the deployment name during command execution.
You can define a deployment that is used as a current one by using the `rasa config use-deployment` command.

```text
Usage:
  rasactl config use-deployment DEPLOYMENT-NAME [flags]
```

```text
Examples:
  # Set the 'example' deployment as the current deployment.
  $ rasactl config use-deployment example
```

```text
Flags:
  -h, --help   help for use-deployment
```

### The `connect rasa` command

Run a local Rasa Open Source server and connect it to a Rasa X deployment.

The command prepares a configuration that's required to connect Rasa X deployment and run a local Rasa server.

It's required to have the 'rasa' command accessible by rasactl.

The command works only if Rasa X deployment runs on a local Kubernetes cluster managed with 'kind'.

```text
Usage:
  rasactl connect rasa [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Connect Rasa Server to Rasa X deployment.
  $ rasactl connect rasa

  # Run a separate rasa server for the Rasa X worker environment.
  $ rasactl connect rasa --run-separate-worker

  # Pass extra arguments to rasa server.
  $ rasactl connect rasa --extra-args="--debug"
```

```text
Flags:
      --extra-args strings    extra arguments for Rasa server
  -h, --help                  help for rasa
  -p, --port int              port to run the Rasa server at (default 5005)
      --run-separate-worker   runs a separate Rasa server for the worker environment
```

### The `auth login` command

Log in to Rasa X / Enterprise.

`auth login` stores credentials in an external credentials store, such as the native keychain of the operating system.

 The following external credential stores will be used:

  *  On macOS: [Apple macOS Keychain Access](https://support.apple.com/en-gb/guide/keychain-access/welcome/mac)
  *  On Linux: [pass](https://www.passwordstore.org/)
  *  On Windows: [Microsoft Windows Credential Manager](https://support.microsoft.com/en-us/windows/accessing-credential-manager-1b5c916a-6a16-889f-8581-fc16e8165ac0)

 You can pass credentials via environment variables:

  *  `RASACTL_AUTH_USER` - username
  *  `RASACTL_AUTH_PASSWORD` - password

 If the environment variables are used, credentials stored in a native keychain are not used.

```text
Usage:
  rasactl auth login [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Login to the 'my-deployment' Rasa X / Enterprise deployment.
  $ rasactl auth login my-deployment

  # Login to Rasa X / Enterprise (login to the currently active deployment).
  $ rasactl auth login

  # Provide a password using STDIN.
  # You can login non-interactively by using the --password-stdin flag to provide a password through STDIN.
  # Using STDIN prevents the password from ending up in the shell’s history.
  $ rasactl auth login --username me --password-stdin
```

```text
Flags:
  -h, --help              help for login
  -p, --password string   password
      --password-stdin    read the password from stdin
  -u, --username string   username
```

***Notice*** For Linux, `pass` is used as credential storage. `pass` must be installed and configured before you use the `rasactl auth` command. Below you can find an example of `pass` installation and configuration.

`pass` installation and configuration for Linux Ubuntu.

1. Install `pass`.

```text
sudo apt-get install pass
```

2. Generate a GPG key.

```text
$  gpg --gen-key
gpg (GnuPG) 2.2.19; Copyright (C) 2019 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

gpg: directory '/home/ubuntu/.gnupg' created
gpg: keybox '/home/ubuntu/.gnupg/pubring.kbx' created
Note: Use "gpg --full-generate-key" for a full featured key generation dialog.

GnuPG needs to construct a user ID to identify your key.

Real name: rasactl
Email address:
You selected this USER-ID:
    "rasactl"

Change (N)ame, (E)mail, or (O)kay/(Q)uit? O
[...]
public and secret key created and signed.
```

3. Init `pass`.

```text
$ pass init rasactl
mkdir: created directory '/home/ubuntu/.password-store/'
Password store initialized for rasactl
```

Now you can use `rasactl auth` on Linux.

```text
$ rasactl ls
CURRENT	NAME             	STATUS 	RASA PRODUCTION	RASA WORKER	ENTERPRISE	VERSION
*       wonderful-gagarin	Running	2.8.1          	2.8.1      	inactive  	0.42.0
$ rasactl auth login
Username: me
Password:
Successfully logged.
```

***Troubleshooting*** If you see `Error: exit status 2: gpg: decryption failed: No secret key` error you should export the following environment variable `export GPG_TTY="$(tty)"`.

### The `auth logout` command

Removes credentials from an external credentials store, such as the native keychain of the operating system.

```text
Usage:
  rasactl auth logout [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Remove access credentials (use the currently active deployment).
  $ rasactl auth logout

  # Remove access credentials for the 'my-deployment' deployment.
  $ rasactl auth logout my-deployment
```

```text
Flags:
  -h, --help   help for logout
```

### The `logs` command

Print the logs for a container in a pod. If the pod has only one container, the container name is optional.

```text
Usage:
  rasactl logs [DEPLOYMENT-NAME] [POD] [flags]
```

```text
Examples:
  # Choose a pod and show logs for it (use the currently active deployment).
  $ rasactl logs

  # Show logs from pod rasa-x (use the currently active deployment).
  $ rasactl logs rasa-x

  # Show logs from pod rasa-x for the 'my-deployment' deployment.
  $ rasactl logs my-deployment rasa-x

  # Display only the most recent 10 lines of output in pod rasa-x
  $ rasactl logs rasa-x --tail=10

  # Return snapshot of previous terminated nginx container logs from pod rasa
  $ rasactl logs -p -c nginx rasa

  # Begin streaming the logs from pod rasa-x
  $ rasactl logs -f rasa-x
```

```text
Flags:
  -c, --container string   a container name
  -f, --follow             specify if the logs should be streamed
  -h, --help               help for logs
  -p, --previous           print the logs for the previous instance of the container in a pod if it exists
      --tail int           lines of recent log file to display. Defaults to -1 showing all log lines (default -1)
```

## Enterprise Management Commands

You can manage an Enterprise license via `rasactl`.

```text
manage Rasa Enterprise

Usage:
  rasactl enterprise [command]

Available Commands:
  activate    activate an Enterprise license
  deactivate  deactivate an Enterprise license
```

### The `enterprise activate` command

Activate an Enterprise license.

```text
Usage:
  rasactl enterprise activate [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Activate an Enterprise license (use the currently active deployment).
  $ rasactl enterprise activate

  # Activate an Enterprise license for the 'my-deployment' deployment.
  $ rasactl enterprise activate my-deployment

  # Provide an Enterprise license using STDIN.
  # You can pass an Enterprise license non-interactively by using the --license-stdin flag to provide a license through STDIN.
  # Using STDIN prevents the license from ending up in the shell’s history.
  $ rasactl enterprise activate --license-stdin
```

```text
Flags:
  -h, --help             help for activate
  -l, --license string   an Enterprise license
      --license-stdin    read an Enterprise license from stdin
```

### The `enterprise deactivate` command

Deactivate an Enterprise license.

```text
Usage:
  rasactl enterprise deactivate [DEPLOYMENT-NAME] [flags]
```

```text
Examples:
  # Deactivate an Enterprise license (use the currently active deployment).
  $ rasactl enterprise deactivate

  # Deactivate an Enterprise license for the 'my-deployment' deployment.
  $ rasactl enterprise deactivate my-deployment
```

```text
Flags:
  -h, --help   help for deactivate
```

## Model Management Commands

You can manage models in Rasa X / Enterprise via `rasactl`. Below is a list of commands that help with managing models:

```text
$ rasactl help model
manage models for Rasa X / Enterprise

Usage:
  rasactl model [command]

Available Commands:
  delete      delete a model from Rasa X / Enterprise
  download    download a model from Rasa X / Enterprise
  list        list models stored in Rasa X / Enterprise
  tag         tag a model in Rasa X / Enterprise
  upload      upload model to Rasa X / Enterprise
```

### The `model delete` command

Delete a model from Rasa X / Enterprise.

```text
Usage:
  rasactl model delete [DEPLOYMENT-NAME] MODEL-NAME [flags]

Aliases:
  delete, del
```

```text
Examples:
  # Delete the 'example-model' model (use the currently active deployment).
  $ rasactl model delete example-model

  # Delete the 'example-model' model for the 'my-deployment' deployment.
  $ rasactl model delete my-deployment example-model
```

```text
Flags:
  -h, --help   help for delete
```

### The `model download` command

Download a model from Rasa X / Enterprise to your local machine.

```text
Usage:
  rasactl model download [DEPLOYMENT-NAME] MODEL-NAME [DESTINATION] [flags]
```

```text
Examples:
  # Download the 'example-model' model (use the currently active deployment).
  # If the destination is not defined, the model will be stored in a current working directory.
  $ rasactl model download example-model

  # Download the 'example-model' model for the 'my-deployment' deployment
  # and store it in the /tmp directory.
  $ rasactl model download my-deployment example-model /tmp/example-model.tar.gz
```

```text
Flags:
  -h, --help   help for download
```

### The `model list` command

List all models stored in Rasa X / Enterprise.

```text
Usage:
  rasactl model list [DEPLOYMENT-NAME] [flags]

Aliases:
  list, ls
```

```text
Examples:
  # List all models (use the currently active deployment).
  $ rasactl model list

  # List all models for the 'my-deployment' deployment.
  $ rasactl model list my-deployment
```

```text
Flags:
  -h, --help   help for list
```

### The `model tag` command

Create a tag and assign it to a given model.

Rasa Enterprise allows multiple versions of an assistant to be run simultaneously and served to different users. By default, two environments are defined:

- production
- worker

 If you want to activate a model you have to tag it as 'production'.

 Learn more [here](https://rasa.com/docs/rasa-x/enterprise/deployment-environments/).

```text
Usage:
  rasactl model tag [DEPLOYMENT-NAME] MODEL-NAME TAG [flags]
```

```text
Examples:
  # Tag the 'my-model' model as 'production' (use the currently active deployment)
  $ rasactl model tag my-model production

  # Tag the 'my-model' with the 'test' tag within the 'my-deployment' deployment.
  $ rasactl model tag my-deployment my-model test
```

```text
Flags:
  -h, --help   help for tag
```

### The `model upload` command

Upload a model to Rasa X / Enterprise.

```text
Usage:
  rasactl model upload [DEPLOYMENT-NAME] MODEL-FILE [flags]

Aliases:
  upload, up
```

```text
Examples:
  # Upload the model.tar.gz model file to Rasa X / Enterprise (use the currently active deployment).
  $ rasactl model upload model.tar.gz

  # Upload the model.tar.gz model file to the 'my-deployment' deployment.
  $ rasactl model upload my-deployment model.tag.gz
```

```text
Flags:
  -h, --help   help for upload
```

### Upload a model to Rasa X

The following example shows how to download an existing model and upload it via `rasactl`.

1. Download a model.

```text
$ curl -L https://github.com/RasaHQ/rasa-x-demo/blob/master/models/model.tar.gz?raw=true --output model.tar.gz
[...]
```

2. Upload the download model to Rasa X.

```text
$ rasactl model upload [deployment name] model.tar.gz

Successfully uploaded.
```

You can use the `rasa model list` command to list all available models, e.g

```text
$ rasactl model list [deployment name]
NAME 	VERSION	COMPATIBLE	TAGS	HASH                            	TRAINED AT
model	2.8.2  	true      	none	093dfaad610d330e5f36e6d7dc104d86	05 Aug 21 13:16 UTC
```

## Examples of usage

### Run Rasa X / Enterprise with a local Rasa Server

It is possible to run a Rasa X / Enterprise deployment with a local rasa server. The following example shows how to connect a local rasa server that is installed in a Python environment to a running Rasa X / Enterprise deployment.

1. Install `rasa` on your local machine. More information on how to install `rasa` you can find in the [docs](https://rasa.com/docs/rasa/installation/).
2. Activate a Python environment with installed `rasa` (this step is optional if you don't use a Python environment).

```bash
$ source .venv/bin/activate
$ rasa --version
Rasa Version      :         2.7.0
Minimum Compatible Version: 2.6.0
Rasa SDK Version  :         2.8.1
Rasa X Version    :         None
Python Version    :         3.7.11
Operating System  :         Darwin-20.5.0-x86_64-i386-64bit
Python Path       :         /repos/rasa/.venv/bin/python3.7
```

3. Connect a local rasa server to a Rasa X / Enterprise deployment.

```bash
$ rasactl connect rasa
●∙∙ Starting Rasa Server
(production-worker) 2021-08-09 15:56:45 INFO     root  - Starting Rasa server on http://localhost:5005
(production-worker) 2021-08-09 15:56:45 INFO     rasa.model  - Loading model models/20210804-105240.tar.gz...
(production-worker) /Users/tczekajlo/repos/rasa/.venv/lib/python3.7/site-packages/rasa/utils/train_utils.py:565: UserWarning: model_confidence is set to `softmax`. It is recommended to try using `model_confidence=linear_norm` to make it easier to tune fallback thresholds.
  category=UserWarning,
2021-08-09 15:56:56 INFO     rasa.core.brokers.pika  - Connecting to RabbitMQ ...
(production-worker) 2021-08-09 15:56:56 INFO     rasa.core.brokers.pika  - RabbitMQ connection to '127.0.0.1' was established.
(production-worker) 2021-08-09 15:56:56 INFO     root  - Rasa server is up and running.
```

4. You can check the status of your deployment and see that Rasa version is the same as the rasa version installed locally.

```bash
$ rasactl status
Name:                   	funny-hopper
Status:                 	Running
URL:                    	http://funny-hopper.rasactl.localhost
Version:                	0.42.0
Enterprise:             	inactive
Rasa production version:	2.7.0
Rasa worker version:    	2.7.0
Project path:           	not defined
```

### Run Rasa X / Enterprise with mounted a local Rasa project

The example shows how to run Rasa X / Enterprise deployment with mounted a local rasa project.

1. Install `rasa` on your local machine. More information on how to install `rasa` you can find in the [docs](https://rasa.com/docs/rasa/installation/).
2. Create a rasa project

```bash
$ rasa init
```

3. Start a new Rasa X / Enterprise deployment.

```bash
$ rasactl start --project
```

(The `rasa start --project` command has to be executed in a directory with rasa project. You can use the `--project-path` flag to pass an absolute path to a rasa project.)

4. Open Rasa X / Enterprise in a web browser.

```bash
$ rasactl open
```

### Upgrade Rasa X / Enterprise version

The following example shows how to upgrade Rasa X / Enterprise version for a deployment that already exists.

1. Create the `values.yaml` file with a specific version.

```yaml
# values.yaml
rasax:
  tag: "0.42.0"
eventService:
  tag: "0.42.0"
dbMigrationService:
  tag: "0.42.0"
```

2. Run upgrade.

```bash
$ rasactl upgrade deployment-name --values-file values.yaml
```

### Deploy Rasa X in one of the public cloud providers

The following example shows how to deploy Rasa X in one of the major cloud providers. In the example, GCP (Google Cloud Platform) is used.

1. Create a VM using a Linux base image. You can find detailed information on how to create a VM [here](https://cloud.google.com/compute/docs/instances/create-start-instance).
2. [Install rasactl](#installation) on the VM
3. Start a new deployment by executing the `rasactl start` command.
4. After several minutes you should see details of your deployment.

```text
$ rasactl start
∙∙∙ Ready!

╭ Rasa X ────────────────────────────────╮
│                                        │
│    URL: http://35.184.183.164:30012    │
│    Rasa X version: 0.42.0              │
│    Rasa X password: rasaxlocal         │
│                                        │
╰────────────────────────────────────────╯
```

***Important!*** The Rasa X / Enterprise deployment will be exposed to the public on one of the service node ports (`30000-30100`). Remember to add a rule to firewall configuration that allows for access to the Rasa X deployment.

## Development

Below you can find a setup required for developing `rasactl` locally.

### How to run it?

1. Install go, e.g. by using brew

```text
$ brew install go
```

2. Compile it

```text
$ make build
```

3. Run it

```text
$ ./dist/rasactl
```

### Run unit tests

```text
make test
```

### Kind cluster for developing purposes

1. Install kind and run it

```text
brew install kind
```

2. Prepare configuration for a kind cluster

```text
$ bash kind/generate-config.sh > config.yaml
```

3. Create a kind cluster

```text
$ kind create cluster --config config.yaml
```

After kind is ready, install ingress-nginx:

```text
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
$ kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission
```

## License

Licensed under the Apache License, Version 2.0.
Copyright 2021 Rasa Technologies GmbH. [Copy of the license](LICENSE.txt).

# rasactl

rasactl deploys Rasa X / Enterprise on your local or remote Kubernetes cluster and manages Rasa X / Enterprise deployments.

## Features

- deploy Rasa X / Enterprise
- upgrade Rasa X / Enterprise
- stop/delete/start a running Rasa X / Enterprise deployment
- connect a local Rasa Server to Rasa X / Enterprise
- mount a local Rasa project to Rasa X / Enterprise deployment

## Prequimentes

- Kubernetes cluster >= 1.14+
- kind (for local mode)

## Installation

Coming soon

## Commands

```text
Available Commands:
  add         add existing Rasa X deployment
  completion  generate the autocompletion script for the specified shell
  connect     connect a component to Rasa X
  delete      delete Rasa X deployment
  help        Help about any command
  list        list deployments
  open        open Rasa X in a web browser
  start       start a Rasa X deployment
  status      show deployment status
  stop        stop Rasa X deployment
  upgrade     upgrade Rasa X deployment
```

### The `start` command

The `start` creates a Rasa X deployment or starts stopped deployment if a given deployment already exists.

```text
This command creates a Rasa X deployment or starts stopped deployment if a given deployment already exists.

If the --project or --project-path is used, a Rasa X deployment will be using a local directory with Rasa project.

If a deployment name is not defined, a random name is generated and used as a deployment name.

Usage:
  rasactl start [DEPLOYMENT NAME] [flags]

Examples:
  # Create a Rasa X deployment.
  $ rasactl start

  # Create a Rasa X deployment with custom configuration, e.g the following configuration changes a Rasa X version.
  # All available values: https://github.com/RasaHQ/rasa-x-helm/blob/main/charts/rasa-x/values.yaml
  $ rasactl start --values-file custom-configuration.yaml

  # Create a Rasa X deployment with a defined password.
  $ rasactl start --rasa-x-password mypassword

  # Create a Rasa X deployment that uses a local Rasa project.
  # The command is executed in a Rasa project directory.
  $ rasactl start --project

  # Create a Rasa X deployment with a defined name.
  $ rasactl start my-deployment

Flags:
  -h, --help                          help for start
  -p, --project                       use the current working directory as a project directory, the flag is ignored if the --project-path flag is used
      --project-path string           absolute path to the project directory directory mounted in kind
      --rasa-x-chart-version string   a helm chart version to use
      --rasa-x-password string        Rasa X password (default "rasaxlocal")
      --rasa-x-password-stdin         read the Rasa X password from stdin
      --rasa-x-release-name string    a helm release name to manage (default "rasa-x")
      --values-file string            absolute path to the values file
      --wait-timeout duration         time to wait for Rasa X to be ready (default 10m0s)

Global Flags:
      --config string       config file (default is $HOME/.rasactl.yaml)
      --debug               enable debug output
      --kubeconfig string   absolute path to the kubeconfig file (default "/Users/tczekajlo/.kube/config")
```

### The `stop` command

The `stop` command stops a running Rasa X / Enterprise deployment.

### The `delete` command

The `delete` command deletes a Rasa X / Enterprise deployment.

### The `list` command list Rasa X / Enterprise deployments

The `list` command list deployments.

### The `connect rasa` command

The connect command connects a component to Rasa X / Enterprise, e.g. you can connect a local Rasa Server to a deployment.

```text
Connect Rasa OSS (Open Source Server) to Rasa X deployment.

The command prepares a configuration that's required to connect Rasa X deployment and run a local Rasa server.

It's required to have the 'rasa' command accessible by rasactl.

The command works only if Rasa X deployment uses a local rasa project.

Usage:
  rasactl connect rasa [DEPLOYMENT NAME] [flags]

Examples:
  # Connect Rasa Server to Rasa X deployment.
  $ rasactl connect rasa

  # Run a saparate rasa server for the Rasa X worker environment.
  $ rasactl connect rasa --run-saparate-worker

  # Pass extra arguments to rasa server.
  $ rasactl connect rasa --extra-args="--debug"

Flags:
      --extra-args strings    extra arguments for Rasa server
  -h, --help                  help for rasa
  -p, --port int              port to run the Rasa server at (default 5005)
      --run-saparate-worker   runs a separate Rasa server for the worker environment

Global Flags:
      --config string       config file (default is $HOME/.rasactl.yaml)
      --debug               enable debug output
      --kubeconfig string   absolute path to the kubeconfig file (default "/Users/tczekajlo/.kube/config")
      --verbose             enable verbose output
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
$ rasactl connect rasa funny-hopper
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
$ rasactl status funny-hopper
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
$ sudo rasactl start --project
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

## Development

Below you can find a setup required for developing `rasactl` locally.

### How to run it?

1. Install go

```
$ brew install go
```

2. Compile it

```
$ go build
```

3. Run it

```
$ ./rasactl
```

4. (optional) Make rasactl global

```
$ sudo cp rasactl /usr/local/bin/
```

### Kind cluster for developing purposes

1. Install kind and run it

```
brew install kind
```

2. Prepare configuration for a kind cluster

```
$ bash kind/generate-config.sh > config.yaml
```

3. Create a kind cluster

```
$ kind create cluster --config config.yaml
```

After kind is ready, install ingress-nginx:

```
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
$ kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission
```

### Deploy Rasa X with mounted a local path

1. Go to a rasa project directory

2. Deploy Rasa X

```
$ sudo ./rasactl start -p
```

## Open Rasa X in a web browser

```
$ ./rasactl open
```

## Deploy Rasa X with mounted a local path and a custom Docker image

1. Create a namespace

```
$ kubectl create ns my-test
```

2. Generate a token

```
$ gcloud auth print-access-token
```

3. Create a secret
```
$ kubectl -n my-test create secret docker-registry gcr --docker-server=eu.gcr.io --docker-username=oauth2accesstoken --docker-password=<token>
```

4. Patch the default service account

```
$ kubectl -n my-test patch serviceaccount default -p '{"imagePullSecrets": [{"name": "gcr"}]}'
```

***Notice*** Token is valid for only one hour, after that time you have to delete the `gcr` secret (`kubectl -n my-test delete secret gcr`) and repeat the 2 and 3 steps.

4. Create a deployment with a custom Docker image

```
$ ./rasactl start my-test -p --values-file testdata/test-image.yaml
```

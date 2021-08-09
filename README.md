# rasactl





## How to run it?

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

## Kind cluster for developing purposes

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

## Deploy Rasa X with mounted a local path

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

## Running rasactl

You can use the `help` command to display description and examples for a specific command, e.g. `rasactl help start`.

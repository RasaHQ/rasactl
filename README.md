# rasaxctl

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
$ ./rasaxctl
```

## Kind cluster for developing purposes

1. Install kind and run it

```
brew install kind
```

2. Create a kind cluster

```
kind create cluster --config kind/cluster-configuration.yaml
```

After kind is ready, install ingress-nginx:

```
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
$ kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission
```

## Deploy Rasa X with mounted a local path

```
$ sudo ./rasaxctl start my-project --project-path /path/to/my/project
```

## Open Rasa X in a web browser

```
$ ./rasaxctl open my-project
```

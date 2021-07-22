package types

type KubernetesBackendType string

const (
	KubernetesBackendLocal  KubernetesBackendType = "local"
	KubernetesBackendRemote KubernetesBackendType = "remote"
)

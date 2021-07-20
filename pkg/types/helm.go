package types

import "time"

type RepositorySpec struct {
	Name string
	URL  string
}

type ConfigurationSpec struct {
	Timeout     time.Duration
	ReleaseName string
	Version     string
}

package types

import "time"

type RepositorySpec struct {
	Name string
	URL  string
}

type HelmConfigurationSpec struct {
	Timeout      time.Duration
	ReleaseName  string
	Version      string
	ReuseValues  bool
	StartProject bool
	Atomic       bool
	Wait         bool
}

/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package types

const (
	// RasaXKubernetesConfigMapName stores the config map name that stores configuration files for Rasa X deployment.
	RasaXKubernetesConfigMapName string = "rasa-x-configuration-files"
)

// HealthEndpointsResponse stores a response for the /api/health endpoint.
type HealthEndpointsResponse struct {
	DatabaseMigration DatabaseMigrationSpec `json:"database_migration"`
	Worker            EnvironmentSpec       `json:"worker"`
	Production        EnvironmentSpec       `json:"production"`
}

// DatabaseMigrationSpec stores specification for a database migration response.
type DatabaseMigrationSpec struct {
	Status            string  `json:"status"`
	ProgressInPercent float32 `json:"progress_in_percent"`
}

// VersionEndpointResponse stores a response from the /api/version endpoint.
type VersionEndpointResponse struct {
	Rasa       RasaSpec `json:"rasa"`
	RasaX      string   `json:"rasa-x"`
	Enterprise bool     `json:"enterprise"`
}

// RasaSpec stores specification for the production and worker version.
type RasaSpec struct {
	Production string `json:"production"`
	Worker     string `json:"worker"`
}

// EnvironmentSpec stores specification for a Rasa X / Enterprise environment.
type EnvironmentSpec struct {
	Version                  string `json:"version"`
	MinimumCompatibleVersion string `json:"minimum_compatible_version"`
	Status                   int    `json:"status"`
}

// EnvironmentConfigurationFile specifies the environment.yaml configuration file for Rasa OSS.
type EnvironmentsConfigurationFile struct {
	Rasa RasaSpecEnvironments `yaml:"rasa"`
}

// RasaSpecEnvironments stores specification for the production and worker environment.
type RasaSpecEnvironments struct {
	Production EnvironmentsConfigurationSpec `yaml:"production"`
	Worker     EnvironmentsConfigurationSpec `yaml:"worker"`
}

// EnvironmentsConfugrationSpec stores specification or a given environment.
type EnvironmentsConfigurationSpec struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

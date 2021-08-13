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
	RasaXKubernetesConfigMapName string = "rasa-x-configuration-files"
)

type HealthEndpointsResponse struct {
	DatabaseMigration DatabaseMigrationSpec `json:"database_migration"`
	Worker            EnvironmentSpec       `json:"worker"`
	Production        EnvironmentSpec       `json:"production"`
}

type DatabaseMigrationSpec struct {
	Status            string  `json:"status"`
	ProgressInPercent float32 `json:"progress_in_percent"`
}

type VersionEndpointResponse struct {
	Rasa       RasaSpec `json:"rasa"`
	RasaX      string   `json:"rasa-x"`
	Enterprise bool     `json:"enterprise"`
}

type RasaSpec struct {
	Production string `json:"production"`
	Worker     string `json:"worker"`
}

type EnvironmentSpec struct {
	Version                  string `json:"version"`
	MinimumCompatibleVersion string `json:"minimum_compatible_version"`
	Status                   int    `json:"status"`
}

type EnvironmentsConfigurationFile struct {
	Rasa RasaSpecEnvironments `yaml:"rasa"`
}

type RasaSpecEnvironments struct {
	Production EnvironmentsConfigurationSpec `yaml:"production"`
	Worker     EnvironmentsConfigurationSpec `yaml:"worker"`
}

type EnvironmentsConfigurationSpec struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type AuthEndpointResponse struct {
	AccessToken string `json:"access_token"`
}

type ModelsListEndpointResponse struct {
	Models []ModelSpec
}

type ModelSpec struct {
	Model        string   `json:"model"`
	Hash         string   `json:"hash"`
	TrainedAt    float64  `json:"trained_at"`
	Version      string   `json:"version"`
	Tags         []string `json:"tags"`
	IsCompatible bool     `json:"is_compatible"`
}

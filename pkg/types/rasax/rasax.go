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

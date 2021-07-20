package types

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

package rasa

type CredentialsFile struct {
	Rasa struct {
		Url string `yaml:"url"`
	} `yaml:"rasa"`
	Rest string `yaml:"rest"`
}

type EndpointsFile struct {
	Models       EndpointModelSpec        `yaml:"models"`
	TrackerStore EndpointTrackerStoreSpec `yaml:"tracker_store"`
	EventBroker  EndpointEventBrokerSpec  `yaml:"event_broker"`
}

type EndpointModelSpec struct {
	Url                  string `yaml:"url"`
	Token                string `yaml:"token"`
	WaitTimeBetweenPulls int    `yaml:"wait_time_between_pulls"`
}

type EndpointTrackerStoreSpec struct {
	Type     string `yaml:"type"`
	Dialect  string `yaml:"dialect"`
	Url      string `yaml:"url"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	LoginDb  string `yaml:"login_db"`
}

type EndpointEventBrokerSpec struct {
	Type     string   `yaml:"type"`
	Url      string   `yaml:"url"`
	Port     int32    `yaml:"port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Queues   []string `yaml:"queues"`
}

/*

   event_broker:
     type: "pika"
     url: "{{ template "rasa-x.rabbitmq.host" . }}"
     username: "{{ .Values.rabbitmq.rabbitmq.username }}"
     password: ${RABBITMQ_PASSWORD}
     port: {{ default 5672 .Values.rabbitmq.service.port }}
     {{ if or (regexMatch ".*(a|rc)[0-9]+" .Values.rasa.version) (regexMatch "2.*[0-9]+-full" .Values.rasa.version) -}}
     queues:
     - ${RABBITMQ_QUEUE}
*/

/*
   tracker_store:
     type: sql
     dialect: "postgresql"
     url: {{ template "rasa-x.psql.host" . }}
     port: {{ template "rasa-x.psql.port" . }}
     username: {{ template "rasa-x.psql.username" . }}
     password: ${DB_PASSWORD}
     db: ${DB_DATABASE}
     {{- if .Values.rasa.useLoginDatabase }}
     login_db: {{ template "rasa-x.psql.database" . }}
     {{- end }}
*/

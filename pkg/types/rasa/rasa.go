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
package rasa

// CredentialsFile defines the credential.yaml file used by Rasa OSS.
type CredentialsFile struct {
	Rasa struct {
		URL string `yaml:"url"`
	} `yaml:"rasa"`
	Rest string `yaml:"rest"`
}

// EndpointsFile defines the endpoints.yaml file used by Rasa OSS.
type EndpointsFile struct {
	Models       EndpointModelSpec        `yaml:"models"`
	TrackerStore EndpointTrackerStoreSpec `yaml:"tracker_store"`
	EventBroker  EndpointEventBrokerSpec  `yaml:"event_broker"`
}

// EndpointModelSpec specifies a configuration for a model server.
type EndpointModelSpec struct {
	URL                  string `yaml:"url"`
	Token                string `yaml:"token"`
	WaitTimeBetweenPulls int    `yaml:"wait_time_between_pulls"`
}

// EndpointTrackerStoreSpec specifies a configuration for Tacker Store.
type EndpointTrackerStoreSpec struct {
	Type     string `yaml:"type"`
	Dialect  string `yaml:"dialect"`
	URL      string `yaml:"url"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	LoginDb  string `yaml:"login_db"`
}

// EndpointEventBrokerSpec stores specification for Event Broker configuration.
type EndpointEventBrokerSpec struct {
	Type     string   `yaml:"type"`
	URL      string   `yaml:"url"`
	Port     int32    `yaml:"port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Queues   []string `yaml:"queues"`
}

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

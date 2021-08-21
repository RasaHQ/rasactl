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

const RasaCtlLocalDomain string = "rasactl.localhost"

type RasaCtlFlags struct {
	StartUpgrade RasaCtlStartUpgradeFlags
	Start        RasaCtlStartFlags
	Delete       RasaCtlDeleteFlags
	Status       RasaCtlStatusFlags
	ConnectRasa  RasaCtlConnectRasaFlags
	Global       RasaCtlGlobalFlags
	Auth         RasaCtlAuthFlags
	Model        RasaCtlModelFlags
	Config       RasaCtlConfigFlags
}

type RasaCtlStartUpgradeFlags struct {
	ValuesFile string
}

type RasaCtlStartFlags struct {
	Create             bool
	ProjectPath        string
	Project            bool
	RasaXPassword      string
	RasaXPasswordStdin bool
	UseEdgeRelease     bool
}

type RasaCtlDeleteFlags struct {
	Force bool
	Prune bool
}

type RasaCtlStatusFlags struct {
	Details bool
}

type RasaCtlConnectRasaFlags struct {
	RunSeparateWorker bool
	Port              int
	ExtraArgs         []string
}

type RasaCtlGlobalFlags struct {
	Debug   bool
	Verbose bool
}

type RasaCtlAuthFlags struct {
	Login struct {
		Username      string
		Password      string
		PasswordStdin bool
	}
}

type RasaCtlModelFlags struct {
	Upload struct {
		File string
	}
	Download struct {
		Name     string
		FilePath string
	}
	Tag struct {
		Name  string
		Model string
	}
	Delete struct {
		Name string
	}
}

type RasaCtlConfigFlags struct {
	CreateFile bool
}

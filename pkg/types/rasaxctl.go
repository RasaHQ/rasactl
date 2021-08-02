package types

type RasaXCtlFlags struct {
	StartUpgrade RasaXCtlStartUpgradeFlags
	Start        RasaXCtlStartFlags
	Delete       RasaXCtlDeleteFlags
	Status       RasaXCtlStatusFlags
	ConnectRasa  RasaXCtlConnectRasaFlags
}

type RasaXCtlStartUpgradeFlags struct {
	ValuesFile string
}

type RasaXCtlStartFlags struct {
	ProjectPath   string
	Project       bool
	RasaXPassword string
}

type RasaXCtlDeleteFlags struct {
	Force bool
	Prune bool
}

type RasaXCtlStatusFlags struct {
	Details bool
}

type RasaXCtlConnectRasaFlags struct {
	RunSeparateWorker bool
	Port              int
	ExtraArgs         []string
}

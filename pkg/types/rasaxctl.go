package types

type RasaXCtlFlags struct {
	StartUpgrade RasaXCtlStartUpgradeFlags
	Start        RasaXCtlStartFlags
	Delete       RasaXCtlDeleteFlags
	Status       RasaXCtlStatusFlags
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

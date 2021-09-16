package status

import (
	"fmt"

	"github.com/RasaHQ/rasactl/pkg/types"
	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

func checkVersionConstrains(version *rtypes.VersionEndpointResponse, flags *types.RasaCtlFlags) {

	if flags.Start.Project || flags.Start.ProjectPath != "" {
		localProjectRasaXVersion(version)
	}
}

func localProjectRasaXVersion(version *rtypes.VersionEndpointResponse) {
	// Check if Rasa X version supports with the mount a local path feature

	if !utils.RasaXVersionConstrains(version.RasaX, ">= 1.0.0") {
		YellowBox(
			"Notice",
			fmt.Sprintf(
				"You're using Rasa X in the %s version, mounting a local rasa project is supported for Rasa X / Enterprise >= 1.0.0",
				version.RasaX,
			),
		)
	}
}

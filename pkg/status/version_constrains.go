package status

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/RasaHQ/rasactl/pkg/types"
	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
)

func checkVersionConstrains(version *rtypes.VersionEndpointResponse, flags *types.RasaCtlFlags) {

	if flags.Start.Project || flags.Start.ProjectPath != "" {
		localProjectRasaXVersion(version)
	}
}

func localProjectRasaXVersion(version *rtypes.VersionEndpointResponse) {
	// Check if Rasa X version supports with the mount a local path feature
	c, err := semver.NewConstraint(">= 1.0.0")
	if err != nil {
		return
	}

	v, err := semver.NewVersion(version.RasaX)
	if err != nil {
		return
	}

	if !c.Check(v) {
		YellowBox(
			"Notice",
			fmt.Sprintf(
				"You're using Rasa X in the %s version, mounting a local rasa project is supported for Rasa X / Enterprise >= 1.0.0",
				version.RasaX,
			),
		)
	}
}

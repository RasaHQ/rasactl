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
package status

import (
	"fmt"
	"strings"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/RasaHQ/rasactl/pkg/types"
	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
	"github.com/RasaHQ/rasactl/pkg/utils"
)

// GreenBox prints a green color box in the terminal.
func GreenBox(tittle string, msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		b := box.New(box.Config{Py: 1, Px: 4, Type: "Round", TitlePos: "Top"})

		b.Config.Color = "Green"

		b.Println(tittle, msg)
	}
}

// RedBox prints a red color box in the terminal.
func RedBox(tittle string, msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		b := box.New(box.Config{Py: 1, Px: 4, Type: "Round", TitlePos: "Top"})

		b.Config.Color = "Red"

		b.Println(tittle, msg)
	}
}

// YellowBox prints a yellow box in the terminal.
func YellowBox(tittle string, msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		b := box.New(box.Config{Py: 1, Px: 4, Type: "Round", TitlePos: "Top"})

		b.Config.Color = "Yellow"

		b.Println(tittle, msg)
	}
}

// PrintRasaXStatus prints a box with details for Rasa X deployment.
func PrintRasaXStatus(version *rtypes.VersionEndpointResponse, url string, flags *types.RasaCtlFlags) {
	if !utils.IsDebugOrVerboseEnabled() {

		msg := []string{fmt.Sprintf("URL: %s", url)}

		if version.Rasa.Production != "0.0.0" && version.Rasa.Worker != "0.0.0" {
			msg = append(msg, fmt.Sprintf("Rasa production version: %s", version.Rasa.Production))
			msg = append(msg, fmt.Sprintf("Rasa worker version: %s", version.Rasa.Worker))
		}

		msg = append(msg,
			fmt.Sprintf("Rasa X version: %s\nRasa X password: %s",
				version.RasaX,
				flags.Start.RasaXPassword,
			),
		)

		GreenBox(
			"Rasa X",
			strings.Join(msg, "\n"),
		)

		// Check the URL
		if !utils.IsURLAccessible(url) {
			YellowBox(
				"Hint",
				fmt.Sprintf("It looks like the %s URL is not accessible, check if all needed firewall rules are in place", url),
			)
		}

		checkVersionConstrains(version, flags)
	}
}

package status

import (
	"fmt"
	"strings"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/spf13/viper"
)

func GreenBox(tittle string, msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		b := box.New(box.Config{Py: 1, Px: 4, Type: "Round", TitlePos: "Top"})

		b.Config.Color = "Green"

		b.Println(tittle, msg)
	}
}

func RedBox(tittle string, msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		b := box.New(box.Config{Py: 1, Px: 4, Type: "Round", TitlePos: "Top"})

		b.Config.Color = "Red"

		b.Println(tittle, msg)
	}
}

func YellowBox(tittle string, msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		b := box.New(box.Config{Py: 1, Px: 4, Type: "Round", TitlePos: "Top"})

		b.Config.Color = "Yellow"

		b.Println(tittle, msg)
	}
}

func PrintRasaXStatus(version *types.VersionEndpointResponse, url string) {
	if !utils.IsDebugOrVerboseEnabled() {

		msg := []string{fmt.Sprintf("URL: %s", url)}

		if version.Rasa.Production != "0.0.0" {
			msg = append(msg, fmt.Sprintf("Rasa production version: %s", version.Rasa.Production))
		}

		msg = append(msg, fmt.Sprintf("Rasa worker version: %s\nRasa X version: %s\nRasa X password: %s", version.Rasa.Worker, version.RasaX, viper.GetString("rasa-x-password")))

		GreenBox(
			"Rasa X",
			strings.Join(msg, "\n"),
		)

		// Check the URL
		if !utils.IsURLAccessible(url) {
			YellowBox(
				"Hint",
				fmt.Sprintf("Looks like the %s URL is not accessible, check if all needed firewall rules are in place", url),
			)
		}
	}
}

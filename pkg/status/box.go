package status

import (
	"fmt"

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

func PrintRasaXStatus(version *types.VersionEndpointResponse, url string) {
	if !utils.IsDebugOrVerboseEnabled() {

		GreenBox(
			"Rasa X",
			fmt.Sprintf("URL: %s\nRasa production version: %s\nRasa worker version: %s\nRasa X version: %s\nRasa X password: %s",
				url, version.Rasa.Production, version.Rasa.Worker, version.RasaX, viper.GetString("rasa-x-password")),
		)
	}
}

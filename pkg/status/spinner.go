package status

import (
	"fmt"
	"os"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/utils"
	"github.com/briandowns/spinner"
)

type SpinnerMessage struct {
	spinner *spinner.Spinner
}

func (s *SpinnerMessage) New() {
	s.spinner = spinner.New(spinner.CharSets[69], 200*time.Millisecond, spinner.WithWriter(os.Stderr))
}

func (s *SpinnerMessage) Message(msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		s.spinner.Suffix = fmt.Sprintf(" %s", msg)
		if !s.spinner.Active() {
			s.spinner.Start()
		} else {
			time.Sleep(1 * time.Second)
			s.spinner.Restart()
		}
	}
}

func (s *SpinnerMessage) Stop() {
	if s.spinner.Active() {
		s.spinner.FinalMSG = "\n"
		s.spinner.Stop()
	}
}

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
	"os"
	"time"

	"github.com/briandowns/spinner"

	"github.com/RasaHQ/rasactl/pkg/utils"
)

// SpinnerMessage defines a spinner object.
type SpinnerMessage struct {
	spinner *spinner.Spinner
}

// NewSpinner creates a new spinner object.
func NewSpinner() *SpinnerMessage {
	s := &SpinnerMessage{}
	s.spinner = spinner.New(spinner.CharSets[69], 200*time.Millisecond, spinner.WithWriter(os.Stderr))
	return s
}

// Message adds a message to the Spinner object.
func (s *SpinnerMessage) Message(msg string) {
	if !utils.IsDebugOrVerboseEnabled() {
		s.spinner.Suffix = fmt.Sprintf(" %s", msg)
		if !s.spinner.Active() {
			s.spinner.Start()
			time.Sleep(200 * time.Millisecond)
		} else {
			time.Sleep(200 * time.Millisecond)
			s.spinner.Restart()
		}
	}
}

func (s *SpinnerMessage) Stop() {
	if s.spinner.Active() {
		s.spinner.Stop()
	}
}

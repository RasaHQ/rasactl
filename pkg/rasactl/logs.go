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
package rasactl

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func (r *RasaCtl) Logs(args []string) error {
	pod := ""

	surveyIconsOpts := survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = ""
		icons.Help.Format = "magenta"
	})

	if args[1] != "" {
		pod = args[1]
	} else {

		podsList, err := r.KubernetesClient.GetPods()
		if err != nil {
			return err
		}

		options := []string{}
		for _, p := range podsList.Items {
			options = append(options, p.Name)
		}

		promptPod := &survey.Select{
			Message: "Choose a pod:",
			Options: options,
		}
		if err := survey.AskOne(promptPod, &pod, surveyIconsOpts); err != nil {

			if errors.Is(err, terminal.InterruptErr) {
				fmt.Println("Interrupted")
				return nil
			}

			return err
		}
	}

	// check if a given pod has more then one container
	podData, err := r.KubernetesClient.GetPod(pod)
	if err != nil {
		return err
	}

	if len(podData.Spec.Containers) > 1 && r.Flags.Logs.Container == "" {
		containers := []string{}

		for _, c := range podData.Spec.Containers {
			containers = append(containers, c.Name)
		}

		promptContainer := &survey.Select{
			Message: "Choose a container:",
			Options: containers,
		}
		if err := survey.AskOne(promptContainer, &r.Flags.Logs.Container, surveyIconsOpts); err != nil {

			if errors.Is(err, terminal.InterruptErr) {
				fmt.Println("Interrupted")
				return nil
			}

			return err
		}
	}

	stream, err := r.KubernetesClient.GetLogs(pod).Stream(context.TODO())
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		buf := make([]byte, 2000)
		numBytes, err := stream.Read(buf)

		if numBytes == 0 {
			break
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		msg := string(buf[:numBytes])
		fmt.Print(msg)
	}

	return nil
}

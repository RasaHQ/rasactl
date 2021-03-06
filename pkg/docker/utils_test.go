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
package docker_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/RasaHQ/rasactl/pkg/docker"
)

var _ = Describe("Utils", func() {

	Context("check Docker version constrains", func() {
		It("Should error", func() {
			err := docker.VersionConstrains("20.9.0")
			Expect(err).To(Not(BeNil()))
		})

		It("Should not error", func() {
			err := docker.VersionConstrains("20.10.1")
			Expect(err).To(BeNil())
		})
	})

	Context("Check skip Docker version check", func() {
		viper.AutomaticEnv() // read in environment variables that match
		viper.SetEnvPrefix("rasactl")

		It("Shoud skip", func() {
			os.Setenv("RASACTL_SKIP_DOCKER_VERSION_CHECK", "true")

			skip := docker.SkipVersionConstrainsCheck()
			Expect(skip).To(BeTrue())
		})

		It("Should not skip when env var sets to false", func() {
			os.Setenv("RASACTL_SKIP_DOCKER_VERSION_CHECK", "false")

			skip := docker.SkipVersionConstrainsCheck()
			Expect(skip).To(BeFalse())
		})

		It("Should not skip when env var does not set", func() {
			skip := docker.SkipVersionConstrainsCheck()
			Expect(skip).To(BeFalse())
		})
	})
})

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
package helm_test

import (
	"os"

	"github.com/RasaHQ/rasactl/pkg/helm"
	"github.com/RasaHQ/rasactl/pkg/logger"
	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	var (
		client helm.Interface
		log    logr.Logger
		flags  *types.RasaCtlFlags
		err    error
	)

	BeforeEach(func() {
		flags = &types.RasaCtlFlags{
			StartUpgrade: types.RasaCtlStartUpgradeFlags{
				ValuesFile: "../../testdata/values.yaml",
			},
			Global: types.RasaCtlGlobalFlags{
				Verbose: false,
				Debug:   false,
			},
		}

		log = logger.New(flags)
	})

	JustBeforeEach(func() {
		client, err = helm.New(
			&helm.Helm{
				Namespace: "test-namespace",
				Flags:     flags,
				Log:       log,
			},
		)
	})

	Describe("Read the values file", func() {
		Context("render template", func() {
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should render template", func() {
				os.Setenv("RASACTL_TEST_VERSION", "0.0.0")
				err := client.ReadValuesFile()
				values := client.GetValues()

				expect := map[string]interface{}{
					"rasax": map[string]interface{}{
						"podLabels": map[string]interface{}{
							"rasactl":       "true",
							"test_template": float64(1),
							"test_version":  "0.0.0",
						},
					},
				}

				Expect(err).NotTo(HaveOccurred())
				Expect(values).Should(Equal(expect))

			})
		})
	})

})

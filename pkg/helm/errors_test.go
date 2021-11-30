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
	"errors"

	"github.com/RasaHQ/rasactl/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {

	Describe("Check errors", func() {
		waitForConditionError := "timed out waiting for the condition. Check your deployment status manually with `rasactl status`"
		Context("timed out waiting for the condition", func() {
			It("should match", func() {
				target := errors.New("timed out waiting for the condition")
				returnError := helm.ErrorTimeoutWaitForCondition(target)
				Expect(returnError).Should(MatchError(waitForConditionError))
			})

			It("should not match", func() {
				target := errors.New("some fake error")
				returnError := helm.ErrorTimeoutWaitForCondition(target)
				Expect(returnError).ShouldNot(MatchError(waitForConditionError))
			})

		})
	})

})

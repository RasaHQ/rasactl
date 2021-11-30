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
package helm

import (
	"errors"
)

// ErrorTimeoutWaitForCondition wraps the "timed out waiting for the condition" error returned by helm.
// If an error is different it returns the original message.
func ErrorTimeoutWaitForCondition(err error) error {
	timeoutWaitForCondition := "timed out waiting for the condition"
	if err.Error() == timeoutWaitForCondition {
		return errors.New("timed out waiting for the condition. Check your deployment status manually with `rasactl status`")
	}
	return err
}

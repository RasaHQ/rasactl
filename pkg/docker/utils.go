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
package docker

import (
	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

const (
	// Env var used for warning on/off
	dockerVersionWarningEnv string = "skip_docker_version_check"
)

// DockerVersionConstrains checks if the Docker CLI version is
// within constrains boundaries.
func VersionConstrains(dockerVersion string) error {
	constraint := ">= 20.10.0"

	if dockerVersion == "" {
		return xerrors.Errorf(
			"Rasactl has an issue reading Docker version. Are you sure Docker is properly installed?",
		)
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return err
	}

	v, err := semver.NewVersion(dockerVersion)
	if err != nil {
		return err
	}

	if !c.Check(v) {
		return xerrors.Errorf(
			"The Docker version is incompatible with rasactl. The version you use is %s"+
				", but rasactl requires Docker %s", dockerVersion, constraint)
	}

	return nil
}

// SkipVersionConstrainsCheck skips Docker version check
// if the `RASACTL_SKIP_DOCKER_VERSION_CHECK` environment variable is set to `true`.
func SkipVersionConstrainsCheck() bool {
	viper.SetDefault(dockerVersionWarningEnv, "false")
	return viper.GetBool(dockerVersionWarningEnv)
}

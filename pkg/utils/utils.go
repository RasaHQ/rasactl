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
package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"github.com/RasaHQ/rasactl/pkg/types"
)

type NetworkError string

// validName is a regular expression for names.
// See: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
var validName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

const (
	NetworkErrorConnectionRefused NetworkError = "ConnectionRefused"
)

func IsDebugOrVerboseEnabled() bool {
	return viper.GetBool("debug") || viper.GetBool("verbose")
}

func CheckNetworkError(err error) (NetworkError, error) {
	var urlError url.Error

	if errors.As(err, urlError) {
		if urlError.Op == "Get" {
			return NetworkErrorConnectionRefused, err
		}
	}

	return "", err
}

// MergeMaps merges maps into one.
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		mergo.Map(&result, m, mergo.WithOverride) //nolint:errcheck
	}
	return result
}

func ValidateName(name string) error {

	if !validName.MatchString(name) {
		return errors.Errorf(
			"Invalid name: \"%s\": a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', "+
				"and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '%s')",
			name, validName.String())
	}

	return nil
}

// IsURLAccessible returns `true` if a client can connect to a given URL.
func IsURLAccessible(address string) bool {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 9,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 3 * time.Second,
			}).Dial,
		},
	}
	req, _ := http.NewRequest("GET", address, nil)
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
}

func readStatusFile(path string, log logr.Logger) (string, error) {
	file := fmt.Sprintf("%s/.rasactl", path)

	log.Info("Reading a status file", "file", file)

	if _, err := os.Stat(file); err != nil {
		log.Info("Status file doesn't exist", "file", file)
		return "", nil
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetActiveNamespace returns an active namespace, it checks the .rasactl file
// if exists and reads a namespace from the file or it gets the current-deployment
// from the rasactl configuration file.
func GetActiveNamespace(log logr.Logger) string {
	log.V(1).Info("Getting active namespace")
	path, err := os.Getwd()
	if err != nil {
		log.V(1).Info("Can't get active namespace", "error", err)
		return ""
	}

	namespace, err := readStatusFile(path, log)
	if err != nil {
		log.V(1).Info("Can't get active namespace", "error", err)
		return ""
	}

	if namespace == "" && viper.GetString("current-deployment") != "" {
		log.V(1).Info("Using namespace name from the rasactl configuration file")
		namespace = viper.GetString("current-deployment")
	}

	return strings.TrimSuffix(namespace, "\n")
}

// AskForConfirmation waits for a input to confirm an operation and returns `true`
// if the input == 'yes'.
func AskForConfirmation(s string, retry int, in io.Reader) (bool, error) {
	r := bufio.NewReader(in)

	for ; retry > 0; retry-- {
		fmt.Printf("%s [yes/no]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			return false, nil
		}

		if len(res) < 2 {
			continue
		}

		response := strings.ToLower(strings.TrimSpace(res))

		if response == "yes" || response == "no" {
			return strings.ToLower(strings.TrimSpace(res)) == "yes", nil
		}
		fmt.Println("You have to put 'yes' or 'no'")
		continue
	}

	return false, nil

}

// Check if a given command exists.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ReadCredentials reads a username and a password from input.
func ReadCredentials(flags *types.RasaCtlFlags) (string, string, error) {
	var username, password string

	if flags.Auth.Login.Username != "" {
		username = flags.Auth.Login.Username
	} else {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Username: ")
		user, err := reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		username = user
	}

	if flags.Auth.Login.Password != "" {
		fmt.Println("WARNING! Using the --password flag is insecure. Use the --password-stdin flag.")
		password = flags.Auth.Login.Password
	} else if flags.Auth.Login.PasswordStdin {
		pass, err := GetPasswordStdin()
		if err != nil {
			return "", "", err
		}
		password = pass
	} else {
		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return "", "", err
		}
		password = string(bytePassword)
	}

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// ReadLicense reads an Enterprise license from input.
func ReadLicense(flags *types.RasaCtlFlags) (string, error) {
	var license string

	if flags.Enterprise.Activate.License != "" {
		fmt.Println("WARNING! Using the --license flag is insecure. Use the --license-stdin flag.")
		license = flags.Enterprise.Activate.License
	} else if flags.Enterprise.Activate.LicenseStdin {
		l, err := GetPasswordStdin()
		if err != nil {
			return "", err
		}
		license = l
	} else {
		fmt.Print("License: ")
		byteLicense, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return "", err
		}
		license = string(byteLicense)
	}

	return strings.TrimSpace(license), nil
}

// GetPasswordStdin reads a password from STDIN.
func GetPasswordStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return line, err
	}
	return strings.TrimSuffix(line, "\n"), nil
}

// HelmChartVersionConstrains checks if the rasa-x-helm chart version
// is within constraints boundaries.
func HelmChartVersionConstrains(helmChartVersion string) error {
	constraint := fmt.Sprintf(">= %s, < 5.0.0", types.HelmChartVersionRasaX)

	if helmChartVersion == "" {
		return nil
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return err
	}

	v, err := semver.NewVersion(helmChartVersion)
	if err != nil {
		return err
	}

	if !c.Check(v) {
		return fmt.Errorf(
			"the helm chart version is incorrect, the version that you want to use is %s"+
				", use the helm chart in version %s", helmChartVersion, constraint)
	}

	return nil
}

// RasaXVersionConstrains checks if Rasa X version is within constraints boundaries.
func RasaXVersionConstrains(version string, constraint string) bool {

	// containers, e.g >= 1.0
	// ref: https://github.com/Masterminds/semver
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return c.Check(v)
}

// StringSliceToJSON converts [][]string{} to JSON.
func StringSliceToJSON(d [][]string) (string, error) {

	data := map[string]interface{}{}

	for _, item := range d {
		field := strings.ToLower(item[0])
		field = strings.ReplaceAll(field, " ", "_")
		field = strings.ReplaceAll(field, ":", "")

		data[field] = item[1]
	}

	jsonByte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonByte), nil
}

// CheckHelmChartDir checks if a local Helm chart directory exists.
func CheckHelmChartDir() {
	name := strings.TrimSpace(types.HelmChartNameRasaX)
	yellowColor := color.New(color.FgYellow)

	if _, err := os.Stat(name); err == nil {
		yellowColor.Printf("WARNING: In your current working directory is located the %s directory, it'll be used as a source for the helm chart.\n", name)
	}
}

// GetRasaXURLEnv returns Rasa X URL passed via environment variables.
func GetRasaXURLEnv(namespace string) string {
	var rasaXURLNamespace string

	rasaXURL := viper.GetString("rasa_x_url")

	if namespace != "" {
		ns := strings.ReplaceAll(namespace, "-", "_")
		rasaXURLNamespace = viper.GetString(fmt.Sprintf("rasa_x_url_%s", ns))
	}

	if rasaXURLNamespace != "" {
		return rasaXURLNamespace
	} else if rasaXURL != "" {
		return rasaXURL
	}

	return ""
}

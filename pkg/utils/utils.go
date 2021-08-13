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

	"github.com/RasaHQ/rasactl/pkg/types"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"
	"golang.org/x/term"
)

type NetworkError string

// validName is a regular expression for names.
// See: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
var validName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

const (
	NetworkErrorConnectionRefused NetworkError = "ConnectionRefused"
)

func IsDebugOrVerboseEnabled() bool {
	if viper.GetBool("debug") || viper.GetBool("verbose") {
		return true
	}

	return false
}

func CheckNetworkError(err error) (NetworkError, error) {
	switch t := err.(type) {
	case *url.Error:
		if t.Op == "Get" {
			return NetworkErrorConnectionRefused, err
		}
	}
	return "", err
}

// MergeMaps merges maps into one.
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		mergo.Map(&result, m, mergo.WithOverride)
	}
	return result
}

// AddHostToEtcHosts adds a host with a given IP address to /etc/hosts,
// for Windows it is C:\Windows\System32\Drivers\etc\hosts.
func AddHostToEtcHosts(host, ip string) error {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	hosts.AddHost(ip, host)

	if err := hosts.Save(); err != nil {
		return errors.Errorf("Can't add a host, try to run the command as administrator, error: %s", err)
	}

	return nil
}

// DeleteHostToEtcHosts removes a host from /etc/hosts (linux, darwin), or
// C:\Windows\System32\Drivers\etc\hosts (Windows).
func DeleteHostToEtcHosts(host string) error {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	hosts.RemoveHost(host)

	if err := hosts.Save(); err != nil {
		return errors.Errorf("Can't remove a host, try to run the command as administrator, error: %s", err)
	}

	return nil
}

func ValidateName(name string) error {

	if !validName.MatchString(name) {
		return errors.Errorf("Invalid name: \"%s\": a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '%s')",
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
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 3 * time.Second,
			}).Dial,
		},
	}
	req, _ := http.NewRequest("GET", address, nil)
	if _, err := client.Do(req); err != nil {
		return false
	}
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

// GetActiveNamespace returns a active namespace, it checks the .rasactl file if exists
// and read a namespace from the file.
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
		} else {
			fmt.Println("You have to put 'yes' or 'no'")
			continue
		}
	}

	return false, nil

}

// Check if a given command exists.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

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
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}
		password = string(bytePassword)
	}

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func GetPasswordStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return line, err
	}
	return strings.TrimSuffix(line, "\n"), nil
}

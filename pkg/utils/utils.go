package utils

import (
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"
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

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {

	result := make(map[string]interface{})
	for _, m := range maps {
		mergo.Map(&result, m, mergo.WithOverride)
	}
	return result
}

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

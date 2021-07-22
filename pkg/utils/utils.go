package utils

import (
	"net/url"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"
)

type NetworkError string

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

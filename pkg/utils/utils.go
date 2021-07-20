package utils

import (
	"net/url"

	"github.com/spf13/viper"
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

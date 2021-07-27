package utils

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

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

// Source: https://cs.opensource.google/go/go/+/master:src/net/ip.go;l=133
// This function can be deleted after Go 1.17 is released.
// IsPrivate reports whether ip is a private address, according to
// RFC 1918 (IPv4 addresses) and RFC 4193 (IPv6 addresses).
func IsPrivate(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		// Following RFC 1918, Section 3. Private Address Space which says:
		//   The Internet Assigned Numbers Authority (IANA) has reserved the
		//   following three blocks of the IP address space for private internets:
		//     10.0.0.0        -   10.255.255.255  (10/8 prefix)
		//     172.16.0.0      -   172.31.255.255  (172.16/12 prefix)
		//     192.168.0.0     -   192.168.255.255 (192.168/16 prefix)
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	// Following RFC 4193, Section 8. IANA Considerations which says:
	//   The IANA has assigned the FC00::/7 prefix to "Unique Local Unicast".
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}

func ValidateName(name string) error {

	if !validName.MatchString(name) {
		return errors.Errorf("Invalid name: \"%s\": a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '%s')",
			name, validName.String())
	}

	return nil
}

func IsURLAccessible(address string) bool {
	ipPort, err := url.Parse(address)
	if err != nil {
		return false
	}
	if _, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ipPort.Host, ipPort.Port())); err != nil {
		return false
	}

	return true
}

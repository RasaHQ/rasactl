package credentials

import (
	"github.com/docker/docker-credential-helpers/credentials"
)

type Helper interface {
	Add(creds *credentials.Credentials) error
	Get(serverURL string) (string, string, error)
	Delete(serverURL string) error
}

type Credentials struct {
	Helper Helper
}

func (c *Credentials) Set(label, url, user, secret string) error {
	cr := &credentials.Credentials{
		ServerURL: url,
		Username:  user,
		Secret:    secret,
	}

	credentials.SetCredsLabel(label)
	return c.Helper.Add(cr)
}

func (c *Credentials) Get(label, url string) (string, string, error) {
	credentials.SetCredsLabel(label)
	return c.Helper.Get(url)
}

func (c *Credentials) Delete(label, url string) error {
	credentials.SetCredsLabel(label)
	return c.Helper.Delete(url)
}

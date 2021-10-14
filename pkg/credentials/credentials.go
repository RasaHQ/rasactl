package credentials

import (
	"fmt"

	"github.com/docker/docker-credential-helpers/credentials"
)

type Helper interface {
	Add(creds *credentials.Credentials) error
	Get(name string) (string, string, error)
	Delete(name string) error
}

type Credentials struct {
	Helper    Helper
	Namespace string
}

func (c *Credentials) Set(name, user, secret string) error {
	cName := fmt.Sprintf("%s-%s", name, c.Namespace)
	cr := &credentials.Credentials{
		ServerURL: fmt.Sprintf("https://%s", cName),
		Username:  user,
		Secret:    secret,
	}

	credentials.SetCredsLabel(cName)
	return c.Helper.Add(cr)
}

func (c *Credentials) Get(name string) (string, string, error) {
	cName := fmt.Sprintf("%s-%s", name, c.Namespace)
	credentials.SetCredsLabel(cName)
	return c.Helper.Get(fmt.Sprintf("https://%s", cName))
}

func (c *Credentials) Delete(name string) error {
	cName := fmt.Sprintf("%s-%s", name, c.Namespace)
	credentials.SetCredsLabel(cName)
	return c.Helper.Delete(fmt.Sprintf("https://%s", cName))
}

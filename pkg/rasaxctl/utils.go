package rasaxctl

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (r *RasaXCTL) writeStatusFile(path string) error {
	d := []byte(r.Namespace)
	file := fmt.Sprintf("%s/.rasaxctl", path)

	r.Log.Info("Writing a status file", "file", file)

	if err := ioutil.WriteFile(file, d, 0644); err != nil {
		return err
	}

	return nil
}

func (r *RasaXCTL) readStatusFile(path string) (string, error) {
	file := fmt.Sprintf("%s/.rasaxctl", path)

	r.Log.Info("Reading a status file", "file", file)

	if _, err := os.Stat(file); err != nil {
		r.Log.Info("Status file doesn't exist", "file", file)
		return "", nil
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *RasaXCTL) GetActiveNamespace() (string, error) {
	r.Log.V(1).Info("Getting active namespace")
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	namespace, err := r.readStatusFile(path)
	if err != nil {
		return "", nil
	}

	return namespace, nil
}

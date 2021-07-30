package rasaxctl

import (
	"fmt"
	"io/ioutil"
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

package gse

import (
	"encoding/json"
	"os/exec"
)

// SetEnabled sets the enabled extensions from a slice.
func SetEnabled(enabled []string) error {
	json, err := json.Marshal(&enabled)
	if err != nil {
		return err
	}

	err = exec.Command("gsettings", "set", "org.gnome.shell", "enabled-extensions", string(json)).Run()
	if err != nil {
		return err
	}

	return nil
}

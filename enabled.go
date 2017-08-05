package gse

import (
	"os/exec"

	yaml "gopkg.in/yaml.v2"
)

// Enabled returns the currently enabled extensions.
func Enabled() ([]string, error) {
	stdout, err := exec.Command("gsettings", "get", "org.gnome.shell", "enabled-extensions").Output()
	if err != nil {
		return nil, err
	}

	if string(stdout) == "@as []\n" {
		return make([]string, 0), nil
	}

	var enabled []string
	// the output of gsettings is not valid JSON, as it uses single quotes, but it is valid YAML
	err = yaml.Unmarshal(stdout, &enabled)
	if err != nil {
		return nil, err
	}

	return enabled, nil
}

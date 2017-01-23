package gse

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os/exec"

	"gopkg.in/yaml.v2"
)

// GNOMEVersion defines a version of the GNOME Shell.
type GNOMEVersion struct {
	Major string `xml:"platform"`
	Minor string `xml:"minor"`
	Patch string `xml:"micro"`
}

const baseURL = "https://extensions.gnome.org"
const bold = csi + "1m"
const csi = "\x1b["
const normal = csi + "0m"

func (version *GNOMEVersion) String() string {
	return fmt.Sprintf("%s.%s.%s", version.Major, version.Minor, version.Patch)
}

func getEnabledExtensions() ([]string, error) {
	stdout, err := exec.Command("gsettings", "get", "org.gnome.shell", "enabled-extensions").Output()

	if err != nil {
		return nil, err
	}

	var enabled []string

	// the output of gsettings is not valid JSON, as it uses single quotes, but it
	// is valid YAML
	err = yaml.Unmarshal(stdout, &enabled)

	if err != nil {
		return nil, err
	}

	return enabled, nil
}

// GetGNOMEVersion returns the current version of GNOME.
func GetGNOMEVersion() (*GNOMEVersion, error) {
	file, err := ioutil.ReadFile("/usr/share/gnome/gnome-version.xml")

	if err != nil {
		return nil, err
	}

	version := new(GNOMEVersion)

	err = xml.Unmarshal(file, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func setEnabledExtensions(enabled []string) error {
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

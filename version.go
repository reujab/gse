package gse

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// GNOMEVersion defines a version of the GNOME Shell.
type GNOMEVersion struct {
	Major string `xml:"platform"`
	Minor string `xml:"minor"`
	Patch string `xml:"micro"`
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

func (version *GNOMEVersion) String() string {
	return fmt.Sprintf("%s.%s.%s", version.Major, version.Minor, version.Patch)
}

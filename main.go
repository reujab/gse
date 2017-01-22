package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type gnomeVersion struct {
	Major string `xml:"platform"`
	Minor string `xml:"minor"`
	Patch string `xml:"micro"`
}

func main() {
	fmt.Println(getGNOMEVersion())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getGNOMEVersion() string {
	file, err := ioutil.ReadFile("/usr/share/gnome/gnome-version.xml")

	check(err)

	data := new(gnomeVersion)

	check(xml.Unmarshal(file, data))

	return fmt.Sprintf("%s.%s.%s", data.Major, data.Minor, data.Patch)
}

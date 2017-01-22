package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"gopkg.in/urfave/cli.v1"
)

type gnomeVersion struct {
	Major string `xml:"platform"`
	Minor string `xml:"minor"`
	Patch string `xml:"micro"`
}

const baseURL = "https://extensions.gnome.org"
const bold = csi + "1m"
const csi = "\x1b["
const normal = csi + "0m"

func main() {
	app := cli.NewApp()

	app.Usage = "A GNOME Shell extension manager"
	app.Commands = []cli.Command{
		{
			Action:    install,
			Name:      "install",
			ShortName: "i",
			Usage:     "Installs an extension by id",
		},
		{
			Action: search,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page, p",
					Value: 1,
				},
			},
			Name:      "search",
			ShortName: "s",
			Usage:     "Searches for an extension",
		},
	}

	check(app.Run(os.Args))
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

func getHomeDir() string {
	usr, err := user.Current()

	check(err)

	return usr.HomeDir
}

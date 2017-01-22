package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/urfave/cli.v1"
)

type extensionQuery struct {
	Extensions []*struct {
		Creator     string `json:"creator"`
		Description string `json:"description"`
		Name        string `json:"name"`
		UUID        string `json:"uuid"`
		ID          int    `json:"pk"`
	} `json:"extensions"`
}

type gnomeVersion struct {
	Major string `xml:"platform"`
	Minor string `xml:"minor"`
	Patch string `xml:"micro"`
}

const baseURL = "https://extensions.gnome.org"

func main() {
	app := cli.NewApp()

	app.Usage = "A GNOME Shell extension manager"
	app.Commands = []cli.Command{
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

func search(ctx *cli.Context) error {
	args := ctx.Args()

	if len(args) > 1 {
		return cli.ShowCommandHelp(ctx, "search")
	}

	query, _ := url.ParseQuery("")

	query.Add("page", ctx.String("page"))
	query.Add("search", args.First())
	query.Add("shell_version", getGNOMEVersion())

	res, err := http.Get(baseURL + "/extension-query/?" + query.Encode())

	check(err)

	if res.StatusCode != 200 {
		return cli.NewExitError("404", 1)
	}

	bytes, err := ioutil.ReadAll(res.Body)

	check(err)

	extensions := new(extensionQuery)

	err = json.Unmarshal(bytes, extensions)

	check(err)

	// TODO: prettier output
	for i, extension := range extensions.Extensions {
		if i != 0 {
			fmt.Println()
		}

		fmt.Printf("%s%s%s - %d - %s\n", "\x1b[1m", extension.Name, "\x1b[0m", extension.ID, extension.Description)
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/urfave/cli.v1"
)

type extensionQuery struct {
	Extensions []*struct {
		Description string `json:"description"`
		Name        string `json:"name"`
		UUID        string `json:"uuid"`
	} `json:"extensions"`
}

func search(ctx *cli.Context) error {
	args := ctx.Args()

	if len(args) > 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
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

	body, err := ioutil.ReadAll(res.Body)

	check(err)

	extensions := new(extensionQuery)

	check(json.Unmarshal(body, extensions))

	// TODO: prettier output
	for i, extension := range extensions.Extensions {
		if i != 0 {
			fmt.Println()
		}

		fmt.Printf("%s%s%s - %s - %s\n", bold, extension.Name, normal, extension.UUID, extension.Description)
	}

	return nil
}

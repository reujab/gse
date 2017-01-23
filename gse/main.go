package main

import (
	"fmt"
	"os"

	"github.com/reujab/gse"

	"gopkg.in/urfave/cli.v1"
)

const bold = csi + "1m"
const csi = "\x1b["
const normal = csi + "0m"

func main() {
	app := cli.NewApp()

	app.Usage = "A GNOME Shell extension manager"
	app.Commands = []cli.Command{
		{
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()

				if len(args) == 0 {
					return cli.ShowCommandHelp(ctx, ctx.Command.Name)
				}

				for _, arg := range args {
					err := gse.Install(arg, true)

					if err != nil {
						return err
					}
				}

				return nil
			},
			Name:      "install",
			ShortName: "i",
			Usage:     "Installs an extension by id",
		},
		{
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()

				if len(args) > 1 {
					return cli.ShowCommandHelp(ctx, ctx.Command.Name)
				}

				version, err := gse.GetGNOMEVersion()

				if err != nil {
					return err
				}

				extensions, err := gse.Search(args.First(), ctx.String("page"), version.String())

				if err != nil {
					return err
				}

				// TODO: prettier output
				for i, extension := range extensions {
					if i != 0 {
						fmt.Println()
					}

					fmt.Printf("%s%s%s - %s - %s\n", bold, extension.Name, normal, extension.UUID, extension.Description)
				}

				return nil
			},
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

	err := app.Run(os.Args)

	if err != nil {
		panic(err)
	}
}

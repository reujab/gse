package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
)

type extensionDetails struct {
	UUID     string `json:"uuid"`
	Versions map[string]struct {
		PK int `json:"pk"`
	} `json:"shell_version_map"`
}

var idRegex = regexp.MustCompile(`^\d+$`)
var installErr = cli.NewExitError("Failed to install some extensions.", 1)
var uuidRegex = regexp.MustCompile(`.+@.+\..+`)

func install(ctx *cli.Context) (exitErr error) {
	args := ctx.Args()

	if len(args) == 0 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	for _, arg := range args {
		query := make(url.Values)

		if idRegex.MatchString(arg) {
			query.Add("pk", arg)
		} else if uuidRegex.MatchString(arg) {
			query.Add("uuid", arg)
		} else {
			exitErr = installErr

			fmt.Fprintln(os.Stderr, bold+arg+normal+" is not an ID or a UUID.")

			continue
		}

		res, err := http.Get(baseURL + "/extension-info/?" + query.Encode())

		check(err)

		if res.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "%d on /extension-info/?%s\n", res.StatusCode, query.Encode())

			exitErr = installErr

			continue
		}

		body, err := ioutil.ReadAll(res.Body)

		check(err)

		details := new(extensionDetails)

		check(json.Unmarshal(body, details))

		var pk int

		for _, version := range details.Versions {
			if version.PK > pk {
				pk = version.PK
			}
		}

		query = make(url.Values)

		query.Add("version_tag", strconv.Itoa(pk))

		res, err = http.Get(baseURL + "/download-extension/" + details.UUID + ".shell-extension.zip?" + query.Encode())

		check(err)

		if res.StatusCode != 200 {
			panic("non-200 status")
		}

		body, err = ioutil.ReadAll(res.Body)

		check(err)

		reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))

		check(err)

		extensionDir := filepath.Join(getHomeDir(), ".local", "share", "gnome-shell", "extensions", details.UUID)

		for _, file := range reader.File {
			info := file.FileInfo()
			dest := filepath.Join(extensionDir, file.Name)

			if info.IsDir() {
				check(os.MkdirAll(dest, info.Mode()))
			} else {
				check(os.MkdirAll(filepath.Dir(dest), 0755))

				openedFile, err := file.Open()

				check(err)

				destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, info.Mode())

				check(err)

				_, err = io.Copy(destFile, openedFile)

				check(err)
			}
		}

		stdout, err := exec.Command("gsettings", "get", "org.gnome.shell", "enabled-extensions").Output()

		check(err)

		var enabled []string

		// the output of gsettings is not valid JSON, as it uses single quotes, but it
		// is valid YAML
		check(yaml.Unmarshal(stdout, &enabled))

		alreadyEnabled := false

		for _, extension := range enabled {
			if extension == details.UUID {
				alreadyEnabled = true

				break
			}
		}

		if !alreadyEnabled {
			enabled = append(enabled, details.UUID)

			json, err := json.Marshal(&enabled)

			check(err)
			check(exec.Command("gsettings", "set", "org.gnome.shell", "enabled-extensions", string(json)).Run())
		}
	}

	return
}

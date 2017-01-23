package gse

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
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v2"
)

var idRegex = regexp.MustCompile(`^\d+$`)
var uuidRegex = regexp.MustCompile(`.+@.+\..+`)

// Install installs an extension by ID (pk) or UUID.
func Install(arg string, enable bool) error {
	query := make(url.Values)

	if idRegex.MatchString(arg) {
		query.Add("pk", arg)
	} else if uuidRegex.MatchString(arg) {
		query.Add("uuid", arg)
	} else {
		return fmt.Errorf("%s%s%s is not an ID or a UUID", bold, arg, normal)
	}

	res, err := http.Get(baseURL + "/extension-info/?" + query.Encode())

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%d on /extension-info/?%s\n", res.StatusCode, query.Encode())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	details := new(Extension)

	err = json.Unmarshal(body, details)

	if err != nil {
		return err
	}

	var pk int

	for _, version := range details.Versions {
		if version.PK > pk {
			pk = version.PK
		}
	}

	query = make(url.Values)

	query.Add("version_tag", strconv.Itoa(pk))

	res, err = http.Get(baseURL + "/download-extension/" + details.UUID + ".shell-extension.zip?" + query.Encode())

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		panic("non-200 status")
	}

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))

	if err != nil {
		return err
	}

	usr, err := user.Current()

	if err != nil {
		return err
	}

	extensionDir := filepath.Join(usr.HomeDir, ".local", "share", "gnome-shell", "extensions", details.UUID)

	for _, file := range reader.File {
		info := file.FileInfo()
		dest := filepath.Join(extensionDir, file.Name)

		if info.IsDir() {
			err = os.MkdirAll(dest, info.Mode())

			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(dest), 0755)

			if err != nil {
				return err
			}

			openedFile, err := file.Open()

			if err != nil {
				return err
			}

			destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, info.Mode())

			if err != nil {
				return err
			}

			_, err = io.Copy(destFile, openedFile)

			if err != nil {
				return err
			}
		}
	}

	if enable {
		stdout, err := exec.Command("gsettings", "get", "org.gnome.shell", "enabled-extensions").Output()

		if err != nil {
			return err
		}

		var enabled []string

		// the output of gsettings is not valid JSON, as it uses single quotes, but it
		// is valid YAML
		err = yaml.Unmarshal(stdout, &enabled)

		if err != nil {
			return err
		}

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

			if err != nil {
				return err
			}

			err = exec.Command("gsettings", "set", "org.gnome.shell", "enabled-extensions", string(json)).Run()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

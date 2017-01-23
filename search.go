package gse

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type extensionQuery struct {
	Extensions []*Extension `json:"extensions"`
}

// Extension is the response of https://extensions.gnome.org/extension-query/.
type Extension struct {
	Creator     string `json:"creator"`
	CreatorURL  string `json:"creator_url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	PK          int    `json:"pk"`
	Screenshot  string `json:"screenshot"`
	UUID        string `json:"uuid"`
	Versions    map[string]struct {
		PK      int `json:"pk"`
		Version int `json:"version"`
	} `json:"shell_version_map"`
}

// Search searches for an extension.
func Search(search, page, version string) ([]*Extension, error) {
	query := make(url.Values)

	query.Add("page", page)
	query.Add("search", search)
	query.Add("shell_version", version)

	res, err := http.Get(baseURL + "/extension-query/?" + query.Encode())

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("404")
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	extensions := new(extensionQuery)

	err = json.Unmarshal(body, extensions)

	if err != nil {
		return nil, err
	}

	return extensions.Extensions, nil
}

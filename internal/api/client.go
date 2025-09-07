package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type Client struct {
	url         string
	Profile     string
	apiKey      string
	ProfileName string
}

func NewClient(url, profile, apiKey string) Client {
	newClient := Client{url, profile, apiKey, ""}
	profileInfo, err := newClient.ProfileInfo()
	if err != nil {
		return newClient
	}
	newClient.ProfileName = profileInfo.Name
	return newClient
}

func (c Client) Request(uri string, params url.Values) ([]byte, error) {
	query := url.Values{
		"from":  {util.ResultWindow},
		"limit": {util.ResultLimit},
	}
	for k, v := range params {
		query.Add(k, v[0])
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		util.Log.Error("error creating request", "error", err)
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	// The profiles endpoint does not accept query parameters.
	if !strings.HasSuffix(uri, fmt.Sprintf("profiles/%s", c.Profile)) {
		req.URL.RawQuery = query.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		util.Log.Error("error reading response body", "error", err)
		return nil, err
	}

	return body, nil
}

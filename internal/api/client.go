package api

import (
	"io"
	"net/http"
	"net/url"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type Client struct {
	url     string
	profile string
	apiKey  string
}

func NewClient(url, profile, apiKey string) Client {
	return Client{url, profile, apiKey}
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
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		util.Log.Error("error reading response body", "error", err)
		return nil, err
	}

	return body, nil
}

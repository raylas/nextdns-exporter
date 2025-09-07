package api

import (
	"encoding/json"
	"fmt"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type ProfileInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type ProfilesResponse struct {
	Profile ProfileInfo `json:"data"`
}

func (c Client) ProfileInfo() (*ProfileInfo, error) {
	profileURL := fmt.Sprintf("%s/profiles/%s", c.url, c.Profile)

	body, err := c.Request(profileURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	profileResponse := ProfilesResponse{}

	err = json.Unmarshal(body, &profileResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	return &profileResponse.Profile, nil
}

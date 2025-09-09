package api

import (
	"encoding/json"
	"fmt"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type StatusResponse struct {
	Statuses []Status `json:"data"`
}

type Status struct {
	Status  string `json:"status"`
	Queries int    `json:"queries"`
}

type StatusMetrics struct {
	TotalQueries   float64
	AllowedQueries float64
	BlockedQueries float64
}

func (c Client) CollectStatus() (*StatusMetrics, error) {
	statusesURL := fmt.Sprintf("%s/profiles/%s/analytics/status", c.url, c.Profile)

	statusResponse := StatusResponse{}
	metrics := StatusMetrics{}

	body, err := c.Request(statusesURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, status := range statusResponse.Statuses {
		switch status.Status {
		case "default":
			metrics.TotalQueries = float64(status.Queries)
		case "allowed":
			metrics.AllowedQueries = float64(status.Queries)
		case "blocked":
			metrics.BlockedQueries = float64(status.Queries)
		}
	}

	return &metrics, nil
}

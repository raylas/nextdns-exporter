package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type DestinationsResponse struct {
	Destinations []Destination `json:"data"`
}

type Destination struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Queries int    `json:"queries"`
}

type DestinationsMetrics struct {
	Destinations []DestinationMetric
}

type DestinationMetric struct {
	Code    string
	Name    string
	Queries float64
}

func (c Client) CollectDestinations() (*DestinationsMetrics, error) {
	destinationsURL := fmt.Sprintf("%s/profiles/%s/analytics/destinations", c.url, c.profile)

	destinationsResponse := DestinationsResponse{}
	metrics := DestinationsMetrics{}

	params := url.Values{
		"type": {"countries"},
	}

	body, err := c.Request(destinationsURL, params)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &destinationsResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, destination := range destinationsResponse.Destinations {
		destination := DestinationMetric{
			Code:    destination.Code,
			Name:    destination.Name,
			Queries: float64(destination.Queries),
		}

		metrics.Destinations = append(metrics.Destinations, destination)
	}

	return &metrics, nil
}

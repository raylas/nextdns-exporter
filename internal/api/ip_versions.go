package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type IPVersionsResponse struct {
	IPVersions []IPVersion `json:"data"`
}

type IPVersion struct {
	Version int `json:"version"`
	Queries int `json:"queries"`
}

type IPVersionsMetrics struct {
	IPVersions []IPVersionMetric
}

type IPVersionMetric struct {
	Version string
	Queries float64
}

func (c Client) CollectIPVersions() (*IPVersionsMetrics, error) {
	ipVersionsURL := fmt.Sprintf("%s/profiles/%s/analytics/ipVersions", c.url, c.profile)

	ipVersionsResponse := IPVersionsResponse{}
	metrics := IPVersionsMetrics{}

	body, err := c.Request(ipVersionsURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &ipVersionsResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, ipVersion := range ipVersionsResponse.IPVersions {
		ipVersion := IPVersionMetric{
			Version: strconv.Itoa(ipVersion.Version),
			Queries: float64(ipVersion.Queries),
		}

		metrics.IPVersions = append(metrics.IPVersions, ipVersion)
	}

	return &metrics, nil
}

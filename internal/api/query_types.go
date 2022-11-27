package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type QueryTypesResponse struct {
	QueryTypes []QueryType `json:"data"`
}

type QueryType struct {
	Type    int    `json:"type"`
	Name    string `json:"name"`
	Queries int    `json:"queries"`
}

type QueryTypesMetrics struct {
	QueryTypes []QueryTypeMetric
}

type QueryTypeMetric struct {
	Type    string
	Name    string
	Queries float64
}

func (c Client) CollectQueryTypes() (*QueryTypesMetrics, error) {
	queryTypesURL := fmt.Sprintf("%s/profiles/%s/analytics/queryTypes", c.url, c.profile)

	queryTypesResponse := QueryTypesResponse{}
	metrics := QueryTypesMetrics{}

	body, err := c.Request(queryTypesURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &queryTypesResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, queryType := range queryTypesResponse.QueryTypes {
		queryType := QueryTypeMetric{
			Type:    strconv.Itoa(queryType.Type),
			Name:    queryType.Name,
			Queries: float64(queryType.Queries),
		}

		metrics.QueryTypes = append(metrics.QueryTypes, queryType)
	}

	return &metrics, nil
}

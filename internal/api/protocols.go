package api

import (
	"encoding/json"
	"fmt"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type ProtocolsResponse struct {
	Protocols []Protocol `json:"data"`
}

type Protocol struct {
	Protocol string `json:"protocol"`
	Queries  int    `json:"queries"`
}

type ProtocolsMetrics struct {
	Protocols []ProtocolMetric
}

type ProtocolMetric struct {
	Protocol string
	Queries  float64
}

func (c Client) CollectProtocols() (*ProtocolsMetrics, error) {
	protocolsURL := fmt.Sprintf("%s/profiles/%s/analytics/protocols", c.url, c.Profile)

	protocolsResponse := ProtocolsResponse{}
	metrics := ProtocolsMetrics{}

	body, err := c.Request(protocolsURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &protocolsResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, protocol := range protocolsResponse.Protocols {
		protocol := ProtocolMetric{
			Protocol: protocol.Protocol,
			Queries:  float64(protocol.Queries),
		}

		metrics.Protocols = append(metrics.Protocols, protocol)
	}

	return &metrics, nil
}

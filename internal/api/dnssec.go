package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type DNSSECResponse struct {
	Data []DNSSEC `json:"data"`
}

type DNSSEC struct {
	Validated bool `json:"validated"`
	Queries   int  `json:"queries"`
}

type DNSSECMetrics struct {
	Data []DNSSECMetric
}

type DNSSECMetric struct {
	Validated string
	Queries   float64
}

func (c Client) CollectDNSSEC() (*DNSSECMetrics, error) {
	dnssecURL := fmt.Sprintf("%s/profiles/%s/analytics/dnssec", c.url, c.profile)

	dnssecResponse := DNSSECResponse{}
	metrics := DNSSECMetrics{}

	body, err := c.Request(dnssecURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &dnssecResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, data := range dnssecResponse.Data {
		data := DNSSECMetric{
			Validated: strconv.FormatBool(data.Validated),
			Queries:   float64(data.Queries),
		}

		metrics.Data = append(metrics.Data, data)
	}

	return &metrics, nil
}

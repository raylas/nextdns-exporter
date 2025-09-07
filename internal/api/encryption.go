package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type EncryptionResponse struct {
	Data []Encryption `json:"data"`
}

type Encryption struct {
	Encrypted bool `json:"encrypted"`
	Queries   int  `json:"queries"`
}

type EncryptionMetrics struct {
	Data []EncryptionMetric
}

type EncryptionMetric struct {
	Encrypted string
	Queries   float64
}

func (c Client) CollectEncryption() (*EncryptionMetrics, error) {
	encryptionURL := fmt.Sprintf("%s/profiles/%s/analytics/encryption", c.url, c.Profile)

	encryptionResponse := EncryptionResponse{}
	metrics := EncryptionMetrics{}

	body, err := c.Request(encryptionURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &encryptionResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, data := range encryptionResponse.Data {
		data := EncryptionMetric{
			Encrypted: strconv.FormatBool(data.Encrypted),
			Queries:   float64(data.Queries),
		}

		metrics.Data = append(metrics.Data, data)
	}

	return &metrics, nil
}

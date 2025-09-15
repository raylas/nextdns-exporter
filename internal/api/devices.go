package api

import (
	"encoding/json"
	"fmt"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type DevicesResponse struct {
	Devices []Device `json:"data"`
}

type Device struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Model   string `json:"model"`
	LocalIP string `json:"localIp"`
	Queries int    `json:"queries"`
}

type DevicesMetrics struct {
	Devices []DeviceMetric
}

type DeviceMetric struct {
	ID      string
	Name    string
	Model   string
	LocalIP string
	Queries float64
}

func (c Client) CollectDevices() (*DevicesMetrics, error) {
	devicesURL := fmt.Sprintf("%s/profiles/%s/analytics/devices", c.url, c.Profile)

	devicesResponse := DevicesResponse{}
	metrics := DevicesMetrics{}

	body, err := c.Request(devicesURL, nil)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &devicesResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, device := range devicesResponse.Devices {
		device := DeviceMetric{
			ID:      device.ID,
			Name:    device.Name,
			Model:   device.Model,
			LocalIP: device.LocalIP,
			Queries: float64(device.Queries),
		}

		metrics.Devices = append(metrics.Devices, device)
	}

	return &metrics, nil
}

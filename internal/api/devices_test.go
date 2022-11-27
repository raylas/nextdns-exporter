package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectDevices(t *testing.T) {
	expected := &DevicesMetrics{
		Devices: []DeviceMetric{
			{
				ID:      "E8TTX",
				Name:    "Gaming PC",
				Model:   "linux",
				LocalIP: "192.168.1.100",
				Queries: 12,
			},
			{
				ID:      "85C3A",
				Name:    "iPhone",
				Model:   "Apple, Inc.",
				LocalIP: "192.168.1.105",
				Queries: 83,
			},
		},
	}

	devices, err := os.ReadFile("../../fixtures/devices.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(devices))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectDevices()
	if err != nil {
		t.Errorf("error collecting devices data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectDestinations(t *testing.T) {
	expected := &DestinationsMetrics{
		Destinations: []DestinationMetric{
			{
				Code:    "US",
				Name:    "United States of America",
				Queries: 137,
			},
			{
				Code:    "DE",
				Name:    "Germany",
				Queries: 2,
			},
			{
				Code:    "FR",
				Name:    "France",
				Queries: 1,
			},
		},
	}

	destinations, err := os.ReadFile("../../fixtures/destinations.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(destinations))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectDestinations()
	if err != nil {
		t.Errorf("error collecting destinations data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

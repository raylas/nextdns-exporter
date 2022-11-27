package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectDomains(t *testing.T) {
	expected := &DomainsMetrics{
		BlockedDomains: []DomainMetric{
			{
				Domain:  "metrics.icloud.com",
				Root:    "icloud.com",
				Queries: 15005,
			},
			{
				Domain:  "app-measurement.com",
				Queries: 3922,
			},
			{
				Domain:  "notify.bugsnag.com",
				Root:    "bugsnag.com",
				Tracker: "bugsnag",
				Queries: 3760,
			},
		},
	}

	domains, err := os.ReadFile("../../fixtures/domains.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(domains))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectDomains()
	if err != nil {
		t.Errorf("error collecting domains data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectDNSSEC(t *testing.T) {
	expected := &DNSSECMetrics{
		Data: []DNSSECMetric{
			{
				Validated: "false",
				Queries:   183,
			},
			{
				Validated: "true",
				Queries:   4,
			},
		},
	}

	ipVersions, err := os.ReadFile("../../fixtures/dnssec.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(ipVersions))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectDNSSEC()
	if err != nil {
		t.Errorf("error collecting DNSSEC data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

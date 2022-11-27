package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectProtocols(t *testing.T) {
	expected := &ProtocolsMetrics{
		Protocols: []ProtocolMetric{
			{
				Protocol: "DNS-over-HTTPS",
				Queries:  17115,
			},
			{
				Protocol: "UDP",
				Queries:  354,
			},
		},
	}

	protocols, err := os.ReadFile("../../fixtures/protocols.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(protocols))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectProtocols()
	if err != nil {
		t.Errorf("error collecting protocols data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

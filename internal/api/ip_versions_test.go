package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectIPVersions(t *testing.T) {
	expected := &IPVersionsMetrics{
		IPVersions: []IPVersionMetric{
			{
				Version: "4",
				Queries: 392,
			},
			{
				Version: "6",
				Queries: 10,
			},
		},
	}

	ipVersions, err := os.ReadFile("../../fixtures/ip_versions.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(ipVersions))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectIPVersions()
	if err != nil {
		t.Errorf("error collecting IP versions data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

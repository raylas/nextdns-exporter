package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectQueryTypes(t *testing.T) {
	expected := &QueryTypesMetrics{
		QueryTypes: []QueryTypeMetric{
			{
				Type:    "1",
				Name:    "A",
				Queries: 207,
			},
			{
				Type:    "28",
				Name:    "AAAA",
				Queries: 199,
			},
			{
				Type:    "65",
				Name:    "HTTPS",
				Queries: 87,
			},
			{
				Type:    "12",
				Name:    "PTR",
				Queries: 11,
			},
			{
				Type:    "16",
				Name:    "TXT",
				Queries: 5,
			},
		},
	}

	queryTypes, err := os.ReadFile("../../fixtures/query_types.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(queryTypes))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectQueryTypes()
	if err != nil {
		t.Errorf("error collecting query types data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectStatus(t *testing.T) {
	expected := &StatusMetrics{
		TotalQueries:   1587523,
		AllowedQueries: 478,
		BlockedQueries: 80343,
	}

	status, err := os.ReadFile("../../fixtures/status.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(status))
	}))
	defer svr.Close()

	c := NewClient(svr.URL)
	res, err := c.CollectStatus("profile", "apikey")
	if err != nil {
		t.Errorf("error collecting status: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

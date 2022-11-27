package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCollectEncryption(t *testing.T) {
	expected := &EncryptionMetrics{
		Data: []EncryptionMetric{
			{
				Encrypted: "true",
				Queries:   48058,
			},
			{
				Encrypted: "false",
				Queries:   1629,
			},
		},
	}

	ipVersions, err := os.ReadFile("../../fixtures/encryption.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(ipVersions))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "profile", "apikey")
	res, err := c.CollectEncryption()
	if err != nil {
		t.Errorf("error collecting encryption data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestProfileInfo(t *testing.T) {
	expected := &ProfileInfo{
		ID:          "4d4270",
		Name:        "Example Profile",
		Fingerprint: "fp6975362082910ed9",
	}

	profile, err := os.ReadFile("../../fixtures/profile.json")
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(profile))
	}))
	defer svr.Close()

	c := NewClient(svr.URL, "4d4270", "apikey")
	res, err := c.ProfileInfo()
	if err != nil {
		t.Errorf("error collecting protocols data: %v", err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

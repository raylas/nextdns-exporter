package util

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "bar")

	if GetEnv("FOO", "baz") != "bar" {
		t.Error("expected to get `bar`")
	}
	if GetEnv("BAR", "baz") != "baz" {
		t.Error("expected to get `baz`")
	}
}

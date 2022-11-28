package util

import (
	"os"
	"testing"
)

func TestDefaultEnv(t *testing.T) {
	os.Setenv("FOO", "bar")

	if DefaultEnv("FOO", "baz") != "bar" {
		t.Error("expected to get `bar`")
	}
	if DefaultEnv("BAR", "baz") != "baz" {
		t.Error("expected to get `baz`")
	}
}

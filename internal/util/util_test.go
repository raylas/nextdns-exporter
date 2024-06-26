package util

import (
	"os"
	"path/filepath"
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

func TestInitSecret(t *testing.T) {
	os.Setenv("BAZ", "1223")
	v, err := initSecret("BAZ")
	if err != nil {
		t.Errorf("expected no error: %s", err)
	}
	if v != "1223" {
		t.Errorf("expected 1223: %s", v)
	}

	f := filepath.Join(t.TempDir(), "boz")
	err = os.WriteFile(f, []byte("345"), 0o755)
	if err != nil {
		t.Errorf("expected no error: %s", err)
	}
	os.Setenv("BOZ_FILE", f)
	v, err = initSecret("BOZ")
	if err != nil {
		t.Errorf("expected no error: %s", err)
	}
	if v != "345" {
		t.Errorf("expected 345: %s", v)
	}
}

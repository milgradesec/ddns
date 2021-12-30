package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	if err := os.Setenv("DDNS_PROVIDER", "Cloudflare"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DDNS_ZONE", "example.com"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DDNS_UPDATE_INTERVAL", "3"); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(); err != nil {
		t.Fatal(err)
	}
}

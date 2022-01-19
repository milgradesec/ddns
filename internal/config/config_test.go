package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("DDNS_PROVIDER", "Cloudflare")

	t.Setenv("DDNS_ZONE", "example.com")

	t.Setenv("DDNS_UPDATE_INTERVAL", "3")

	if _, err := Load(); err != nil {
		t.Fatal(err)
	}
}

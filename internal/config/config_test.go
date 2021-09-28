package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	_, err := New("../../test/config.json")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadFromEnv(t *testing.T) {
	if err := os.Setenv("DDNS_PROVIDER", "Cloudflare"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DDNS_UPDATE_INTERVAL", "3"); err != nil {
		t.Fatal(err)
	}

	cfg := &Configuration{}
	if err := cfg.LoadFromEnv(); err != nil {
		t.Fatal(err)
	}
}

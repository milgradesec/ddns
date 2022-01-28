package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("DDNS_PROVIDER", "Cloudflare")
	t.Setenv("DDNS_ZONE", "example.com")
	t.Setenv("DDNS_UPDATE_INTERVAL", "3")
	t.Setenv("DDNS_EXCLUDE", "a.example.com,c.example.com")

	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		domain string
		want   bool
	}{
		{"a.example.com", true},
		{"b.example.com", false},
		{"c.example.com", true},
	}
	for _, tt := range tests {
		if got := config.IsExcluded(tt.domain); got != tt.want {
			t.Errorf("Configuration.IsExcluded(%s) = %v, want %v", tt.domain, got, tt.want)
		}
	}
}

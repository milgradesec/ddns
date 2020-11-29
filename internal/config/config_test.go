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
	err := os.Setenv("CLOUDFLARE_API_KEY", "difiowehfhsahsdshndjqwh")
	if err != nil {
		t.Fatal(err)
	}

	_, err = New("../../test/envconfig.json")
	if err != nil {
		t.Fatal(err)
	}
}

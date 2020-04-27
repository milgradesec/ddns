package config

import "testing"

func TestLoad(t *testing.T) {
	_, err := Load("../../test/cfg.json")
	if err != nil {
		t.Fatal(err)
	}
}

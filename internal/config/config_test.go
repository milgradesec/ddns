package config

import "testing"

func TestLoad(t *testing.T) {
	_, err := Load("../../test/config.json")
	if err != nil {
		t.Fatal(err)
	}
}

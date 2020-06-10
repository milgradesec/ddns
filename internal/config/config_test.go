package config

import "testing"

func TestLoad(t *testing.T) {
	_, err := New("../../test/config.json")
	if err != nil {
		t.Fatal(err)
	}
}

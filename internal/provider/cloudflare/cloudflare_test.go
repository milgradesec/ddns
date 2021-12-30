package cloudflare

import (
	"os"
	"testing"

	"github.com/milgradesec/ddns/internal/config"
)

func TestNew(t *testing.T) {
	if err := os.Setenv("CLOUDFLARE_ZONE", "example.com"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("CLOUDFLARE_API_TOKEN", "XXXxxXXXXXXxXXXXXXxxXXX"); err != nil {
		t.Fatal(err)
	}

	if _, err := New(&config.Configuration{}); err != nil {
		t.Fatal(err)
	}
}

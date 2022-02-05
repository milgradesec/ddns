package cloudflare

import (
	"testing"

	"github.com/milgradesec/ddns/internal/config"
)

func TestNew(t *testing.T) {
	t.Setenv("CLOUDFLARE_API_TOKEN", "XXXxxXXXXXXxXXXXXXxxXXX")

	if _, err := New(&config.Configuration{
		Zone: "example.com",
	}); err != nil {
		t.Fatal(err)
	}
}

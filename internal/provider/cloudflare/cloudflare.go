package cloudflare

import (
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/internal/ip"
)

// API implements provider.API interface
type API struct {
	cfg config.Config

	api *cloudflare.API
	id  string
}

// New creates a Cloudflare DNS provider
func New(cfg config.Config) (*API, error) {
	if cfg.IsEmpty() {
		cfg.APIKey = os.Getenv("CF_API_KEY")
		cfg.Email = os.Getenv("CF_API_EMAIL")
		cfg.Zone = os.Getenv("CF_ZONE_NAME")
	}

	api, err := cloudflare.New(cfg.APIKey, cfg.Email)
	if err != nil {
		return nil, err
	}

	cf := &API{
		api: api,
		cfg: cfg,
	}
	return cf, nil
}

// Name implements Provider interface
func (cf *API) Name() string {
	return "Cloudflare"
}

// UpdateZone implements ProviderAPI interface
func (cf *API) UpdateZone() error {
	if cf.id == "" {
		id, err := cf.api.ZoneIDByName(cf.cfg.Zone)
		if err != nil {
			return err
		}
		cf.id = id
	}

	publicIP, err := ip.GetIP()
	if err != nil {
		return err
	}

	records, err := cf.api.DNSRecords(cf.id, cloudflare.DNSRecord{})
	if err != nil {
		return err
	}

	var update bool
	for _, r := range records {
		switch r.Type {
		case "A":
			if r.Content != publicIP {
				update = true
			}
		}
		if update {
			rr := cloudflare.DNSRecord{
				Type:    r.Type,
				Name:    r.Name,
				Content: publicIP,
				Proxied: r.Proxied,
			}
			if err := cf.api.UpdateDNSRecord(cf.id, r.ID, rr); err != nil {
				return err
			}
			fmt.Printf("%s updated from %s to %s\n", r.Name, r.Content, publicIP)
			update = false
		}
	}
	return nil
}

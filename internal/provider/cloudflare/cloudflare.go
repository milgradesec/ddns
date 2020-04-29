package cloudflare

import (
	"errors"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/pkg/ip"
)

// API implements provider.API interface
type API struct {
	cfg config.Config

	api *cloudflare.API
	id  string
}

// New creates a Cloudflare DNS provider
func New(cfg config.Config) (*API, error) {
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
			return errors.New("cloudflare API error: ZoneIDByName command failed")
		}
		cf.id = id
	}

	publicIP, err := ip.GetIP()
	if err != nil {
		return err
	}

	records, err := cf.api.DNSRecords(cf.id, cloudflare.DNSRecord{})
	if err != nil {
		return errors.New("cloudflare API error: failed to list DNS records")
	}

	for _, r := range records {
		if cf.cfg.IsExcluded(r.Name) {
			continue
		}

		switch r.Type {
		case "A":
			if r.Content != publicIP {
				rr := cloudflare.DNSRecord{
					Type:    r.Type,
					Name:    r.Name,
					Content: publicIP,
					Proxied: r.Proxied,
				}
				if err := cf.api.UpdateDNSRecord(cf.id, r.ID, rr); err != nil {
					return fmt.Errorf("error updating %s: %v", r.Name, err)
				}
				log.Printf("%s updated from %s to %s\n", r.Name, r.Content, publicIP)
			}
		}
	}
	return nil
}

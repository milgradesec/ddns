package cloudflare

import (
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/milgradesec/ddns/pkg/ip"
)

// API implements ProviderAPI interface
type API struct {
	api *cloudflare.API
	id  string
}

// New creates a Cloudflare DNS provider
func New() (*API, error) {
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		return nil, err
	}

	cf := &API{
		api: api,
	}
	return cf, nil
}

// UpdateZone implements ProviderAPI interface
func (cf *API) UpdateZone() error {
	if cf.id == "" {
		id, err := cf.api.ZoneIDByName(os.Getenv("CF_ZONE_NAME"))
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

package cloudflare

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/pkg/ip"

	httpc "github.com/milgradesec/go-libs/http"
	log "github.com/sirupsen/logrus"
)

// API implements provider.API interface.
type API struct {
	cfg *config.Configuration

	cf *cloudflare.API
	id string
}

// New creates a Cloudflare DNS provider.
func New(cfg *config.Configuration) (*API, error) {
	if cfg.GetAuthType() == config.APIKey {
		return newWithAPIKey(cfg)
	}
	return newWithAPIToken(cfg)
}

// Creates a Clouflare Provider using API Key.
func newWithAPIKey(cfg *config.Configuration) (*API, error) {
	api, err := cloudflare.New(cfg.APIKey, cfg.Email, cloudflare.HTTPClient(httpc.NewHTTPClient()))
	if err != nil {
		return nil, err
	}

	cf := &API{
		cf:  api,
		cfg: cfg,
	}
	return cf, nil
}

// Creates a Cloudflare Provider using API Token.
func newWithAPIToken(cfg *config.Configuration) (*API, error) {
	api, err := cloudflare.NewWithAPIToken(cfg.APIToken, cloudflare.HTTPClient(httpc.NewHTTPClient()))
	if err != nil {
		return nil, err
	}

	cf := &API{
		cf:  api,
		cfg: cfg,
	}
	return cf, nil
}

// Name implements the provider.API interface.
func (api *API) Name() string {
	return "Cloudflare"
}

// UpdateZone implements the provider.API interface.
func (api *API) UpdateZone(ctx context.Context) error {
	if api.id == "" {
		id, err := api.cf.ZoneIDByName(api.cfg.Zone)
		if err != nil {
			return fmt.Errorf("cloudflare api error: failed to retrieve zone id: %w", err)
		}
		api.id = id
	}

	publicIP, err := ip.GetIP(ctx)
	if err != nil {
		return err
	}

	records, err := api.cf.DNSRecords(ctx, api.id, cloudflare.DNSRecord{})
	if err != nil {
		return fmt.Errorf("cloudflare api error: failed to list dns records: %w", err)
	}

	for i := range records {
		r := records[i]

		if api.cfg.IsExcluded(r.Name) {
			continue
		}

		if r.Type == "A" {
			if r.Content != publicIP {
				rr := cloudflare.DNSRecord{
					Type:    r.Type,
					Name:    r.Name,
					Content: publicIP,
					Proxied: r.Proxied,
				}

				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()

				if err := api.cf.UpdateDNSRecord(ctx, api.id, r.ID, rr); err != nil {
					return fmt.Errorf("error updating %s: %w", r.Name, err)
				}
				log.Infof("updated %s from %s to %s", r.Name, r.Content, publicIP)
			}
		}
	}
	return nil
}

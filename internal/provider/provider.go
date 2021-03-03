package provider

import "context"

// DNSProvider interface represents an DNSProvider from a dns provider
// able to query and modify records for a domain.
type DNSProvider interface {
	Name() string
	UpdateZone(ctx context.Context) error
}

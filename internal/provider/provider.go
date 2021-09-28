package provider

import "context"

// DNSProvider interface represents a dns provider able
// to query and modify records for a domain.
type DNSProvider interface {
	Name() string
	GetZoneName() string
	UpdateZone(ctx context.Context) error
}

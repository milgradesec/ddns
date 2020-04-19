package provider

// DNSProvider interface represents an API from a dns provider
// able to query and modify records for a domain
type DNSProvider interface {
	UpdateZone() error
}

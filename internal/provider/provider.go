package provider

// ProviderAPI interface represents an API from a dns provider
// able to query and modify records for a domain
type ProviderAPI interface {
	Name() string
	UpdateZone() error
}

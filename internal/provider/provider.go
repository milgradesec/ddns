package provider

// API interface represents an API from a dns provider
// able to query and modify records for a domain.
type API interface {
	Name() string
	UpdateZone() error
}

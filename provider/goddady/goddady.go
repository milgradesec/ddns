package goddady

// API implements ProviderAPI interface
type API struct {
	domain string // GO_DOMAIN_NAME
	apiKey string // GO_API_KEY
	secret string // GO_SECRET
}

// New creates a Goddady DNS provider
func New() *API {
	return nil
}

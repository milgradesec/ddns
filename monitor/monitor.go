package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/milgradesec/ddns/provider"
)

// Monitor detects public IP changes for a dns Zone
type Monitor struct {
	zone     string
	provider provider.ProviderAPI
}

// New creates a new Monitor for a Zone with the selected Provider
func New(zone string, provider provider.ProviderAPI) *Monitor {
	return &Monitor{
		zone:     zone,
		provider: provider,
	}
}

// Run starts monitoring
func (m *Monitor) Run() {
	ticker := time.NewTicker(defaultInterval)
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := m.provider.UpdateZone(); err != nil {
					fmt.Printf("error updating zone %s: %v\n", m.zone, err)
				}
			case <-sighup:
				fmt.Printf("SIGHUP received: updating records for %s\n", m.zone)
				if err := m.provider.UpdateZone(); err != nil {
					fmt.Printf("error updating zone %s: %v\n", m.zone, err)
				}
			}
		}
	}()

	if err := m.provider.UpdateZone(); err != nil {
		fmt.Printf("error updating zone %s: %v\n", m.zone, err)
	}
	<-stop
}

const defaultInterval = 3 * time.Minute

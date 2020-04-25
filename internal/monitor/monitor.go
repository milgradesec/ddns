package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/milgradesec/ddns/internal/provider"
)

const defaultInterval = 3 * time.Minute

// Monitor runs in a infinite loop and triggers provider zone updates
// every 3 min interval, can be triggered at any time by sending a
// SIGHUP signal
type Monitor struct {
	Zone string
	prov provider.ProviderAPI
}

// New creates a new Monitor for a Zone with the selected Provider
func New(zone string, provider provider.ProviderAPI) *Monitor {
	return &Monitor{
		Zone: zone,
		prov: provider,
	}
}

func (m *Monitor) callProvider() {
	if err := m.prov.UpdateZone(); err != nil {
		fmt.Printf("error updating zone %s: %v\n", m.Zone, err)
	}
}

// Run starts monitoring
func (m *Monitor) Run() {
	ticker := time.NewTicker(defaultInterval)
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	stop := make(chan bool)
	go func() {
		m.callProvider()
		for {
			select {
			case <-ticker.C:
				m.callProvider()

			case <-sighup:
				fmt.Printf("SIGHUP received: updating records for %s\n", m.Zone)
				m.callProvider()
			}
		}
	}()
	<-stop
}

package ip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
	httpc "github.com/milgradesec/go-libs/http"
)

const (
	DefaultProviderName = "ipify"
)

type response struct {
	IP string
}

// GetIP returns the current public IP obtained from ipify.org.
func GetIP(ctx context.Context) (string, error) {
	client := httpc.NewHTTPClient()
	b := backoff.NewExponentialBackOff()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org/?format=json", nil)
	if err != nil {
		return "", err
	}

	var lastError error
	for i := 0; i < 5; i++ {
		resp, err := client.Do(req)
		if err != nil {
			lastError = err
			time.Sleep(b.NextBackOff())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", errors.New("invalid status code from ipify.org: " + strconv.Itoa(resp.StatusCode))
		}

		var msg response
		err = json.NewDecoder(resp.Body).Decode(&msg)
		if err != nil {
			return "", fmt.Errorf("failed to decode JSON-encoded response: %w", err)
		}

		ip := net.ParseIP(msg.IP)
		if ip == nil {
			return "", errors.New("invalid ip: " + msg.IP)
		}
		return msg.IP, nil
	}
	return "", fmt.Errorf("failed to reach ipify.org: %w", lastError)
}

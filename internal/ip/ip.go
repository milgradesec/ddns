package ip

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
)

type ipifyResponse struct {
	IP string
}

// GetIP returns the current public IP obtained from ipify.org.
func GetIP() (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	b := backoff.NewExponentialBackOff()

	req, err := http.NewRequest(http.MethodGet, "https://api.ipify.org/?format=json", nil)
	if err != nil {
		return "", err
	}

	for i := 0; i < 5; i++ {
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(b.NextBackOff())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", errors.New("invalid status code from ipify.org: " + strconv.Itoa(resp.StatusCode))
		}

		var msg ipifyResponse
		err = json.NewDecoder(resp.Body).Decode(&msg)
		if err != nil {
			return "", err
		}

		ip := net.ParseIP(msg.IP)
		if ip == nil {
			return "", errors.New("failed to parse ip: " + msg.IP)
		}
		return msg.IP, nil
	}
	return "", errors.New("failed to reach ipify")
}

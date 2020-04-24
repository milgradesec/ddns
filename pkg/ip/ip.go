package ip

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
)

type jsonResponse struct {
	IP string
}

// GetIP returns the current public IP obtained from ipify.org
func GetIP() (string, error) {
	b := backoff.NewExponentialBackOff()
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

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

		var msg jsonResponse
		if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			return "", err
		}

		if resp.StatusCode != 200 {
			return "", errors.New("invalid status code from ipify: " + strconv.Itoa(resp.StatusCode))
		}
		return msg.IP, nil
	}
	return "", errors.New("failed to reach ipify")
}

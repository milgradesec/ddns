package monitor

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
)

// GetIP returns the current public IP obtained from ipify.org
func GetIP() (string, error) {
	b := backoff.NewExponentialBackOff()
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.ipify.org/", nil)
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

		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		if resp.StatusCode != 200 {
			return "", errors.New("invalid status code from ipify: " + strconv.Itoa(resp.StatusCode))
		}

		return string(ip), nil
	}
	return "", errors.New("failed to reach ipify")
}

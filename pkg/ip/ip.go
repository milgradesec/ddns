package ip

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
)

type response struct {
	IP string
}

// GetIP returns the current public IP obtained from ipify.org.
func GetIP() (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			},
		},
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

		var msg response
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

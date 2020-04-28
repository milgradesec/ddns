package ip

import (
	"net"
	"testing"
)

func TestGetIP(t *testing.T) {
	rawIP, err := GetIP()
	if err != nil {
		t.Error(err)
	}

	if ip := net.ParseIP(rawIP); ip == nil {
		t.Errorf("failed to parse ip %s", rawIP)
	}
}

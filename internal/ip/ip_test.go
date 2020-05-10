package ip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	_, err := GetIP()
	if err != nil {
		t.Error(err)
	}
}

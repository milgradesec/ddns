package ip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	if _, err := GetIP(); err != nil {
		t.Error(err)
	}
}

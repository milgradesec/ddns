package ip

import (
	"context"
	"testing"
)

func TestGetIP(t *testing.T) {
	if _, err := GetIP(context.TODO()); err != nil {
		t.Error(err)
	}
}

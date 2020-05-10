package updater

import "testing"

func TestUpdate(t *testing.T) {
	if err := Update("v1.3.0"); err != nil {
		t.Fatal(err)
	}
}

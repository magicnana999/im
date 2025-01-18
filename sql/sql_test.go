package sql

import "testing"

func TestSelect(t *testing.T) {
	err := Select()
	if err != nil {
		t.Error(err)
	}
}

package sql

import "testing"

func TestSelect(t *testing.T) {
	_, e := Select()
	if e != nil {
		t.Error(e)
	}
}

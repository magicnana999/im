package sql

import "testing"

func TestSelect(t *testing.T) {
	err := Select("select user_id from im_user")
	if err != nil {
		t.Error(err)
	}
}

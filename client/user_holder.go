package main

import "sync"

type UserHolder struct {
	m sync.Map
}

func NewUserHolder() *UserHolder {
	return &UserHolder{}
}

func (uh *UserHolder) Load(userID int64) (*User, bool) {
	v, ok := uh.m.Load(userID)
	if ok {
		return v.(*User), true
	} else {
		return nil, false
	}
}

func (uh *UserHolder) Store(user *User) {
	uh.m.Store(user.UserID, user)
}

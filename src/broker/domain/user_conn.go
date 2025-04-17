package domain

import (
	"fmt"
	"github.com/magicnana999/im/define"
	"io"
	"strconv"
	"sync/atomic"
)

// UserConn represents a user's connection to the IM broker.
// It stores connection metadata, status, and I/O interfaces for communication.
type UserConn struct {
	Fd            int           `json:"fd"`            // Fd is the file descriptor of the connection.
	AppId         string        `json:"appId"`         // AppId identifies the application associated with the connection.
	UserId        int64         `json:"userId"`        // UserId is the unique identifier of the user.
	ClientAddr    string        `json:"clientAddr"`    // ClientAddr is the client's network address.
	BrokerAddr    string        `json:"brokerAddr"`    // BrokerAddr is the broker's network address.
	OS            define.OSType `json:"os"`            // OS indicates the client's operating system type.
	ConnectTime   int64         `json:"connectTime"`   // ConnectTime is the Unix timestamp when the connection was established.
	IsLogin       atomic.Bool   `json:"isLogin"`       // IsLogin indicates whether the user is logged in (thread-safe).
	IsClosed      atomic.Bool   `json:"isClosed"`      // IsClosed indicates whether the connection is closed (thread-safe).
	LastHeartbeat atomic.Int64  `json:"lastHeartbeat"` // LastHeartbeat is the Unix timestamp of the last heartbeat (thread-safe).
	Reader        io.Reader     `json:"-"`             // Reader is the input stream for reading data from the connection.
	Writer        io.Writer     `json:"-"`             // Writer is the output stream for writing data to the connection.
}

// Desc returns a string description of the UserConn.
// It combines the client address and the connection label.
func (u *UserConn) Desc() string {
	return fmt.Sprintf("%s#%s", u.ClientAddr, u.Label())
}

// Label returns a unique label for the UserConn.
// The label is formatted as "AppId#UserId#DeviceType" if all components are valid.
// It returns an empty string if AppId, UserId, or DeviceType is invalid.
func (u *UserConn) Label() string {
	if len(u.AppId) == 0 {
		return ""
	}
	if u.UserId == 0 {
		return ""
	}

	dt := u.OS.GetDeviceType()

	if len(dt.String()) == 0 {
		return ""
	}

	return fmt.Sprintf("%s#%s#%s", u.AppId, strconv.FormatInt(u.UserId, 10), dt.String())
}

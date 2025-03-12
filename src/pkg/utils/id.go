package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/rs/xid"
	"sync"
)

var (
	snowFlakeNode *snowflake.Node
	mutex         *sync.Mutex
)

func GenerateXId() string {
	guid := xid.New()
	return guid.String()
}

package id

import (
	"github.com/bwmarrin/snowflake"
	"github.com/rs/xid"
	"sync"
)

var (
	node  *snowflake.Node
	mutex *sync.Mutex
)

func init() {
	snowflake.Epoch = 1577836800000
	n, _ := snowflake.NewNode(0)
	node = n
}
func GenerateXId() string {
	guid := xid.New()
	return guid.String()
}

func SnowflakeID() int64 {
	return int64(node.Generate())
}

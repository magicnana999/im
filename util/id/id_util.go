package id

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/magicnana999/im/util"
	"github.com/rs/xid"
	"github.com/sony/sonyflake"
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

func GetSnowFlakeNode() (*snowflake.Node, error) {
	if snowFlakeNode == nil {
		mutex.Lock()
		defer mutex.Unlock()
		snowFlakeNode, err := snowflake.NewNode(12)
		return snowFlakeNode, err
	}

	if snowFlakeNode == nil {
		fmt.Errorf("could not generate snowflake node")
	}
	return snowFlakeNode, nil
}

func GenerateSonyFlakeId() (uint64, error) {
	var st sonyflake.Settings
	st.MachineID = util.GetRandomUint16
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}

	return sf.NextID()
}

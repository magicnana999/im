package main

import (
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/timandy/routine"
)

func main() {

	gnet.NewClient()
	server := &gnet.Server{
		// 设置自定义的Codec
		Codec: &MyCodec{},
	}

	logger.InfoF("Start im broker %d", routine.Goid())
	server.Start()

}

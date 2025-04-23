package vo

import "github.com/magicnana999/im/api/kitex_gen/api"

type DeliverFail struct {
	M      *api.Message
	UserId int64
	Label  []string
}

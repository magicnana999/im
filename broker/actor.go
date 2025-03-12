package broker

import (
	"github.com/asynkron/protoactor-go/actor"
)

type ActorReceiver interface {
	Receive(ctx actor.Context)
}
type brokerActor struct {
	Handler ActorReceiver
}

func InitBrokerActor(r ActorReceiver) *brokerActor {
	return &brokerActor{Handler: r}
}

func (s *brokerActor) Receive(ctx actor.Context) {
	s.Handler.Receive(ctx)
}

//func (state *brokerActor) Receive(ctx actor.Context) {
//	switch msg := ctx.Message().(type) {
//	case *shared.ChatMessage:
//		fmt.Printf("Broker received message from %s: %s\n", msg.UserId, msg.Content)
//		reply := &shared.ChatResponse{Reply: "Broker handled message"}
//		ctx.Respond(reply)
//	}
//}

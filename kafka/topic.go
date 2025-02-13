package kafka

var (
	Route   = TopicInfo{"msg-route", "msg-route-group"}
	Store   = TopicInfo{"msg-store", "msg-store-group"}
	Offline = TopicInfo{"msg-offline", "msg-offline-group"}
	Push    = TopicInfo{"msg-push", "msg-push-group"}
	Deliver = TopicInfo{"msg-deliver", "msg-deliver-group"}
)

type TopicInfo struct {
	Topic string
	Group string
}

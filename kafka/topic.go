package kafka

var (
	Route    = TopicInfo{"msg-route", "msg-route-group"}
	RouteDLQ = TopicInfo{"msg-route-dlq", "msg-route-dlq-group"}
	Store    = TopicInfo{"msg-store", "msg-store-group"}
	Offline  = TopicInfo{"msg-offline", "msg-offline-group"}
	Push     = TopicInfo{"msg-push", "msg-push-group"}
)

type TopicInfo struct {
	Topic string
	Group string
}

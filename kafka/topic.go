package kafka

var (
	Route = TopicInfo{"im-message-route", "im-message-route-group"}
)

type TopicInfo struct {
	topic string
	group string
}

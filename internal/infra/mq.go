package infra

type QueueEngine interface {
	PublishMessage(topic string, message []byte) error
	Close() error
}

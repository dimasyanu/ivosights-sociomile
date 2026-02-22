package infra

type QueueClient interface {
	PublishMessage(topic string, message []byte) error
	GetPublishedMessages() [][]byte
	Close() error
}

package infra

type QueueClient interface {
	Publish(queue string, message []byte) error
	Consume(queue string) (<-chan []byte, error)
	GetPublishedMessages() [][]byte
	Close() error
}

type ConversationCreatedMessage struct {
	TenantID uint   `json:"tenant_id"`
	CustID   string `json:"cust_id"`
	ConvID   string `json:"conv_id"`
	Message  string `json:"message"`
}

type ConversationStatusUpdatedMessage struct {
	ConvID string `json:"conv_id"`
	Status string `json:"status"`
}

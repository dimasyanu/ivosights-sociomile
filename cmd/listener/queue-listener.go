package listener

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/service"
)

type QueueListener struct {
	mq       infra.QueueClient
	queue    string
	convSvc  *service.ConversationService
	messages <-chan []byte
}

func NewQueueListener(c *config.RabbitMQConfig, db *sql.DB) (*QueueListener, error) {
	mq, err := infra.NewRabbitMQClient(c)
	if err != nil {
		return nil, err
	}
	return &QueueListener{
		mq:      mq,
		queue:   c.Queue,
		convSvc: service.NewConversationService(mysqlrepo.NewConversationRepository(db), mq),
	}, nil
}

func (l *QueueListener) Start(done *sync.WaitGroup) {
	defer done.Done()

	// Listen to the queue and process messages
	var err error
	l.messages, err = l.mq.Consume(l.queue)
	if err != nil {
		return
	}
	for msg := range l.messages {
		log.Printf("Received message: %s", string(msg))
		data := &infra.ConversationCreatedMessage{}
		if err := json.Unmarshal(msg, data); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}
		err := l.convSvc.AssignConversation(data.ConvID)
		if err != nil {
			log.Printf("Failed to assign conversation: %v", err)
		}
	}
}

func (l *QueueListener) Close() {
	l.mq.Close()
}

package listener

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
)

type QueueListener struct {
	mq       infra.QueueClient
	queue    string
	convSvc  *service.ConversationService
	userSvc  *service.UserService
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
		userSvc: service.NewUserService(mysqlrepo.NewUserRepository(db)),
	}, nil
}

func (l *QueueListener) processMessage(msg []byte) {
	// Unmarshal the message
	log.Printf("Received message: %s", string(msg))
	data := &infra.ConversationCreatedMessage{}
	if err := json.Unmarshal(msg, data); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// Process the message
	agent, err := l.userSvc.GetAvailableAgent()
	if err != nil {
		log.Printf("Failed to get available agent: %v", err)
		return
	}
	err = l.convSvc.AssignConversation(data.ConvID, agent.ID)
	if err != nil {
		log.Printf("Failed to assign conversation: %v", err)
	}
}

func (l *QueueListener) Start(done *sync.WaitGroup) {
	defer done.Done()

	// Listen to the queue and process messages
	var err error
	l.messages, err = l.mq.Consume(l.queue)
	if err != nil {
		return
	}

	// Process messages until the channel is closed
	for msg := range l.messages {
		l.processMessage(msg)
	}
}

func (l *QueueListener) Close() {
	l.mq.Close()
}

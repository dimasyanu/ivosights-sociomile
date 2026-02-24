package infra

import (
	"testing"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/stretchr/testify/suite"
)

var envPath = "../../.env"

type RabbitMqTestSuite struct {
	config *config.RabbitMQConfig
	client QueueClient

	suite.Suite
}

func TestRabbitMQTestSuite(t *testing.T) {
	suite.Run(t, new(RabbitMqTestSuite))
}

func (s *RabbitMqTestSuite) SetupSuite() {
	s.config = config.NewRabbitMQConfig(envPath)

	var err error
	s.client, err = NewRabbitMQClient(s.config)
	s.NoError(err)
	s.NotNil(s.client)
}

func (s *RabbitMqTestSuite) TearDownSuite() {
	s.client.Clear()
	s.client.Close()
}

func (s *RabbitMqTestSuite) TestNewRabbitMQConfig() {
	s.NotEmpty(s.config.URL)
	s.NotEmpty(s.config.Exchange)
	s.NotEmpty(s.config.Queue)
	s.NotEmpty(s.config.RoutingKey)
	s.Equal("application/json", s.config.ContentType)
}

func (s *RabbitMqTestSuite) TestNewRabbitMQClient() {
	client, err := NewRabbitMQClient(s.config)
	s.NoError(err)
	s.NotNil(client)
	defer client.Close()
}

func (s *RabbitMqTestSuite) TestPublishMessage() {
	client, err := NewRabbitMQClient(s.config)
	s.NoError(err)
	s.NotNil(client)
	defer client.Close()

	err = client.Publish(ConversationQueue, []byte(`{"message": "Hello, RabbitMQ!"}`))
	s.NoError(err)

	published := client.GetPublishedMessages()
	s.Len(published, 1)
	s.Equal([]byte(`{"message": "Hello, RabbitMQ!"}`), published[0])
}

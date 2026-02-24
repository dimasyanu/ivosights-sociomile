package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	"github.com/google/uuid"
)

func (s *FullCycleTestSuite) webhookReceivedNewMessage() {
	// Use existing tenant ID (assuming the created tenant has ID 1) and a random customer ID
	var tenantID uint
	custID := uuid.New()
	s.db.QueryRow("SELECT id FROM tenants LIMIT 1").Scan(&tenantID)

	// Set up payload for webhook received new message
	msg := "Hello, I need help with my order."
	payload := &models.ChannelPayload{
		TenantID:   tenantID,
		CustomerID: custID,
		Message:    msg,
		SenderType: constant.SenderTypeCustomer,
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Simulate webhook received new message from channel
	req := httptest.NewRequest("POST", "/api/v1/channel/webhook", jsonReader)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(201, res.StatusCode)

	// Check inside database
	convID := uuid.Nil
	convAgentID := uuid.Nil
	s.db.
		QueryRow("SELECT id, assigned_agent_id FROM conversations WHERE tenant_id = ? AND customer_id = UUID_TO_BIN(?)", payload.TenantID, payload.CustomerID).
		Scan(&convID, &convAgentID)
	s.NotEqual(uuid.Nil, convID)
	s.Equal(uuid.Nil, convAgentID) // Initially, no agent should be assigned

	message := ""
	msgConvID := uuid.Nil
	s.db.
		QueryRow("SELECT conversation_id, message FROM messages WHERE conversation_id = UUID_TO_BIN(?)", convID).
		Scan(&msgConvID, &message)
	s.Equal(msg, message)
	s.Equal(convID, msgConvID)

	// Wait for the worker to process the message
	time.Sleep(time.Millisecond * 60)

	// Check if the worker assigned an agent to the conversation
	s.db.
		QueryRow("SELECT id, assigned_agent_id FROM conversations WHERE tenant_id = ? AND customer_id = UUID_TO_BIN(?)", payload.TenantID, payload.CustomerID).
		Scan(&convID, &convAgentID)

	// Check if an agent is assigned to the conversation
	agentID := uuid.Nil
	s.db.
		QueryRow("SELECT id FROM users WHERE email = ?", agentEmail).
		Scan(&agentID)
	s.NotEqual(uuid.Nil, agentID)
	s.Equal(agentID, convAgentID)
}

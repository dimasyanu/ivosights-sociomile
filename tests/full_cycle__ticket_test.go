package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *FullCycleTestSuite) escalateConversationToTicket(convID uuid.UUID) {
	payload := &models.EscalateToTicketRequest{
		Reason: "Need further assistance",
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	req := httptest.NewRequest("POST", "/conversations/"+convID.String()+"/escalate", jsonReader)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)
}

package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (s *FullCycleTestSuite) escalateConversationToTicket(convID uuid.UUID) uuid.UUID {
	payload := &models.EscalateToTicketRequest{
		Title:       "Need further assistance",
		Description: "Customer requires more help after initial response.",
		Priority:    1,
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	req := httptest.NewRequest("POST", "/api/v1/backoffice/conversations/"+convID.String()+"/escalate", jsonReader)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)

	// Check if ticket was created
	ticketID := uuid.Nil
	var title, description, status string
	err = s.db.QueryRow("SELECT id, title, description, status FROM tickets WHERE conversation_id = UUID_TO_BIN(?)", convID.String()).Scan(&ticketID, &title, &description, &status)
	s.Require().NoError(err)
	s.Equal(payload.Title, title)
	s.Equal(payload.Description, description)
	s.Equal(constant.TicketStatusOpen, status)

	return ticketID
}

func (s *FullCycleTestSuite) updateTicketStatus(id uuid.UUID, newStatus string, shouldSucceed bool) {
	payload := &models.UpdateTicketStatusRequest{
		Status: newStatus,
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send update ticket status request
	req := httptest.NewRequest("PATCH", "/api/v1/backoffice/tickets/"+id.String()+"/status", jsonReader)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)
	s.Require().NoError(err)

	if shouldSucceed {
		s.Equal(fiber.StatusOK, resp.StatusCode)
	} else {
		s.Equal(fiber.StatusForbidden, resp.StatusCode)
	}

	// Verify the ticket status in the database
	var status string
	err = s.db.QueryRow("SELECT status FROM tickets WHERE id = UUID_TO_BIN(?)", id.String()).Scan(&status)
	s.Require().NoError(err)

	if shouldSucceed {
		s.Equal(newStatus, status)
		return
	}

	s.NotEqual(newStatus, status)
}

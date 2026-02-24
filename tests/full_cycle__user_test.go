package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
)

func (s *FullCycleTestSuite) createAgentUser() {
	// Set up payload for create user request
	payload := &models.UserCreateRequest{
		Name:           "Test Agent",
		Email:          agentEmail,
		Roles:          []string{"agent"},
		Password:       agentPassword,
		RepeatPassword: agentPassword,
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send create user request with authorization header
	req := httptest.NewRequest("POST", "/api/v1/backoffice/users", jsonReader)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(201, res.StatusCode)
}

package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
)

func (s *FullCycleTestSuite) createTenant() {
	// Set up payload for create tenant request
	payload := &models.TenantCreateRequest{
		Name: "Test Tenant",
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send create tenant request with authorization header
	req := httptest.NewRequest("POST", "/api/v1/backoffice/tenants", jsonReader)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(201, res.StatusCode)

	// Check inside database
	name := ""
	s.db.QueryRow("SELECT name FROM tenants WHERE name = ?", payload.Name).Scan(&name)
	s.Equal("Test Tenant", name)
}

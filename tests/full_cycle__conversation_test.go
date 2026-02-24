package tests

import (
	"encoding/json"
	"net/http/httptest"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/google/uuid"
)

func (s *FullCycleTestSuite) getConversationList() []domain.Conversation {
	// Create and send create tenant request with authorization header
	req := httptest.NewRequest("GET", "/api/v1/backoffice/conversations", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(200, res.StatusCode)

	// Parse response body
	var resBody models.Res[domain.Paginated[domain.Conversation]]
	err = json.NewDecoder(res.Body).Decode(&resBody)
	s.Require().NoError(err)
	s.Len(resBody.Data.Items, 1)
	s.Equal(int64(1), resBody.Data.Total)

	return resBody.Data.Items
}

func (s *FullCycleTestSuite) getConversationDetails(convID uuid.UUID) {
	// Create and send create tenant request with authorization header
	req := httptest.NewRequest("GET", "/api/v1/backoffice/conversations/"+convID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(200, res.StatusCode)

	// Parse response body
	var resBody models.Res[domain.ConversationDetail]
	err = json.NewDecoder(res.Body).Decode(&resBody)
	s.Require().NoError(err)
	s.Equal(convID, resBody.Data.ID)
	s.Equal("Test Tenant", *resBody.Data.TenantName)
}

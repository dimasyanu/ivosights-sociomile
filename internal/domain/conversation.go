package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConversationEntity struct {
	ID              uuid.UUID  `sql:"column:id;primaryKey;autoIncrement"`
	TenantID        uint       `sql:"column:tenant_id;type:int;not null"`
	CustomerID      uuid.UUID  `sql:"column:customer_id;type:binary(16);not null"`
	AssignedAgentID *uuid.UUID `sql:"column:assigned_agent_id;type:binary(16);nullable"`
	Status          string     `sql:"column:status;type:varchar(20);not null"`
	CreatedAt       time.Time  `sql:"column:created_at;autoCreateTime;not null"`
	DeletedAt       *time.Time `sql:"column:deleted_at;nullable"`

	TenantName *string `sql:"-"`
	AgentName  *string `sql:"-"`
	AgentEmail *string `sql:"-"`
}

type Conversation struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uint       `json:"tenant_id"`
	CustomerID      uuid.UUID  `json:"customer_id"`
	AssignedAgentID *uuid.UUID `json:"assigned_agent_id"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ConversationDetail struct {
	AssignedAgentName  *string `json:"assigned_agent_name"`
	AssignedAgentEmail *string `json:"assigned_agent_email"`
	TenantName         *string `json:"tenant_name"`

	Conversation
}

type ConversationFilter struct {
	TenantID        *uint      `form:"tenant_id"`
	CustomerID      *uuid.UUID `form:"customer_id"`
	AssignedAgentID *uuid.UUID `form:"assigned_agent_id"`
	Status          *string    `form:"status"`
	Filter
}

func (e *ConversationEntity) ToDto() *Conversation {
	return &Conversation{
		ID:              e.ID,
		TenantID:        e.TenantID,
		CustomerID:      e.CustomerID,
		AssignedAgentID: e.AssignedAgentID,
		Status:          e.Status,
		CreatedAt:       e.CreatedAt,
	}
}

func (e *ConversationEntity) ToDetailDto() *ConversationDetail {
	return &ConversationDetail{
		AssignedAgentName:  e.AgentName,
		AssignedAgentEmail: e.AgentEmail,
		TenantName:         e.TenantName,
		Conversation: Conversation{
			ID:              e.ID,
			TenantID:        e.TenantID,
			CustomerID:      e.CustomerID,
			AssignedAgentID: e.AssignedAgentID,
			Status:          e.Status,
			CreatedAt:       e.CreatedAt,
		},
	}
}

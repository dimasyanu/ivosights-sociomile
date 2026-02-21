package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConversationEntity struct {
	ID              int        `sql:"column:id;primaryKey;autoIncrement"`
	TenantID        int        `sql:"column:tenant_id;type:int;not null"`
	CustomerID      uuid.UUID  `sql:"column:customer_id;type:binary(16);not null"`
	AssignedAgentID *uuid.UUID `sql:"column:assigned_agent_id;type:binary(16);nullable"`
	Status          string     `sql:"column:status;type:varchar(20);not null"`
	CreatedAt       time.Time  `sql:"column:created_at;autoCreateTime;not null"`
	DeletedAt       *time.Time `sql:"column:deleted_at;nullable"`
}

type Conversation struct {
	ID              int        `json:"id"`
	TenantID        int        `json:"tenant_id"`
	CustomerID      uuid.UUID  `json:"customer_id"`
	AssignedAgentID *uuid.UUID `json:"assigned_agent_id"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

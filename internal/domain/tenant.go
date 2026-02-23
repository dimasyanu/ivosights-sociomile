package domain

import "time"

type TenantEntity struct {
	ID        uint       `sql:"column:id;primaryKey;autoIncrement"`
	Name      string     `sql:"column:name;type:varchar(255);not null"`
	CreatedAt time.Time  `sql:"column:created_at;autoCreateTime;not null"`
	UpdatedAt time.Time  `sql:"column:updated_at;autoUpdateTime;not null"`
	DeletedAt *time.Time `sql:"column:deleted_at;index"`
}

type Tenant struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *TenantEntity) ToDto() *Tenant {
	return &Tenant{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

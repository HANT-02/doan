package base_struct

import (
	"doan/pkg/utils"
	"time"
)

type BaseDomainEntity struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewBaseDomainEntity() BaseDomainEntity {
	id := utils.GenerateUUID()
	now := time.Now()
	return BaseDomainEntity{
		ID:        id,
		CreatedAt: now,
	}
}

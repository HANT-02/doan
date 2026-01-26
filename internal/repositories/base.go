package repositories

import (
	"context"
)

type BaseRepository[T any] interface {
	GetTable() string
	GetByCondition(ctx context.Context, condition *CommonCondition) (*Pagination[T], error)
	GetTotal(ctx context.Context, condition *CommonCondition) (uint64, error)
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, id interface{}, updatedData map[string]interface{}) error
	UpdateWithIDs(ctx context.Context, ids []string, updatedData map[string]interface{}) error
	SoftDelete(ctx context.Context, id interface{}) error
	HardDelete(ctx context.Context, ids interface{}) error
	GetByID(ctx context.Context, id interface{}) (*T, error)
}

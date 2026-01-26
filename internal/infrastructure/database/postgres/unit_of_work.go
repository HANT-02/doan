package postgres

import (
	"context"
	repository "doan/internal/repositories"
	apperrors "doan/pkg/error"
	"errors"

	"gorm.io/gorm"
)

const (
	ContextKeyDBTransaction = "gorm_transaction"
)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) repository.UnitOfWork {
	return &unitOfWork{db: db}
}

func (uow *unitOfWork) Begin(ctx context.Context) (context.Context, error) {
	if ctx.Err() != nil {
		return ctx, ctx.Err()
	}

	// Check if transaction already exists
	if _, ok := ctx.Value(ContextKeyDBTransaction).(*gorm.DB); ok {
		return ctx, nil
	}

	tx := uow.db.Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}
	return context.WithValue(ctx, ContextKeyDBTransaction, tx), nil
}

func (uow *unitOfWork) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(ContextKeyDBTransaction).(*gorm.DB)
	if !ok {
		return apperrors.NewInternalError("transaction not found in context", errors.New("transaction not found"))
	}
	return tx.Commit().Error
}

func (uow *unitOfWork) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(ContextKeyDBTransaction).(*gorm.DB)
	if !ok {
		return apperrors.NewInternalError("transaction not found in context", errors.New("transaction not found"))
	}
	return tx.Rollback().Error
}

func GetDb(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(ContextKeyDBTransaction).(*gorm.DB)
	if !ok || tx == nil {
		return db.WithContext(ctx)
	}
	return tx.WithContext(ctx)
}

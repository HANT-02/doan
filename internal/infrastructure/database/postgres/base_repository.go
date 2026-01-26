package postgres

import (
	"context"
	repository "doan/internal/repositories"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	apperrors "doan/pkg/error"
	"doan/pkg/logger"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type baseRepository[T any] struct {
	base_struct.BaseDependency
	db    *gorm.DB
	table string
}

func NewBaseRepository[T any](
	log logger.Logger,
	configManager config.Manager,
	db *gorm.DB,
	table string,
) repository.BaseRepository[T] {
	return &baseRepository[T]{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: configManager,
		},
		db:    db,
		table: table,
	}
}
func (r *baseRepository[T]) GetTable() string {
	return r.table
}
func (r *baseRepository[T]) GetTotal(ctx context.Context, condition *repository.CommonCondition) (uint64, error) {
	op := "baseRepository.GetTotal"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	// Xây dựng query có áp dụng filter
	conditionWithoutPaging := repository.CommonCondition{}
	if condition != nil {
		conditionWithoutPaging.Conditions = condition.Conditions
	}
	var err error
	db, err = BuildQuery(db, &conditionWithoutPaging)
	if err != nil {
		xErr := apperrors.Wrap(err, op, "Failed while building query")
		r.Log.Error(ctx, "Failed while building query", "error", xErr)
		return 0, xErr
	}

	var total int64
	if err := db.Table(r.table).Where("deleted_at is null").Count(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		if errors.Is(err, gorm.ErrModelValueRequired) {
			return uint64(total), nil
		}
		xErr := apperrors.Wrap(err, op, "Failed while counting total")
		r.Log.Error(ctx, "Failed while counting total", "error", xErr)
		return 0, xErr
	}

	return uint64(total), nil
}
func (r *baseRepository[T]) GetByCondition(ctx context.Context, condition *repository.CommonCondition) (*repository.Pagination[T], error) {
	op := "baseRepository.GetByCondition"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	total, err := r.GetTotal(ctx, condition)
	if err != nil {
		xErr := apperrors.Wrap(err, op, "Failed to get total")
		r.Log.Error(ctx, "Failed to get total", "error", xErr)
		return nil, xErr
	}
	totalPage := uint64(1)
	itemsPerPage := uint64(0)
	currentPage := uint64(1)
	if total > 0 && condition.Paging.Limit > 0 {
		totalPage = total / condition.Paging.Limit
		if total%condition.Paging.Limit > 0 {
			totalPage++
		}
		itemsPerPage = condition.Paging.Limit
		currentPage = condition.Paging.Page
	} else {
		totalPage = 1
		itemsPerPage = total
		currentPage = 1
	}

	db, err = BuildQuery(db, condition)
	if err != nil {
		xErr := apperrors.Wrap(err, op, "Failed while building query")
		r.Log.Error(ctx, "Failed while building query", "error", xErr)
		return nil, xErr
	}
	var results []*T
	if condition.Columns != nil && len(condition.Columns) > 0 {
		db = db.Select(condition.Columns)
	}
	if err := db.Table(r.table).Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if len(results) == 0 {
			return nil, nil
		}
		return &repository.Pagination[T]{
			Data: results,
			Meta: repository.Meta{
				ItemsPerPage: itemsPerPage,
				TotalItems:   total,
				CurrentPage:  currentPage,
				TotalPages:   totalPage,
			},
		}, nil
	}
	return &repository.Pagination[T]{
		Data: results,
		Meta: repository.Meta{
			ItemsPerPage: itemsPerPage,
			TotalItems:   total,
			CurrentPage:  currentPage,
			TotalPages:   totalPage,
		},
	}, nil
}
func (r *baseRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	op := "baseRepository.Create"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	if err := db.Table(r.table).Create(entity).Error; err != nil {
		xErr := apperrors.Wrap(err, op, "Failed to create entity")
		r.Log.Error(ctx, "Failed to create entity", "error", xErr)
		return nil, xErr
	}

	return entity, nil
}
func (r *baseRepository[T]) Update(ctx context.Context, id interface{}, updatedData map[string]interface{}) error {
	op := "baseRepository.Update"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	// GORM: Update theo primary key
	updatedData["updated_at"] = gorm.Expr("CURRENT_TIMESTAMP") // Cập nhật thời gian cập nhật
	result := db.Table(r.table).Where("id = ?", id).Updates(updatedData)
	if result.Error != nil {
		xErr := apperrors.Wrap(result.Error, op, "Failed to update entity")
		r.Log.Error(ctx, "Failed to update entity", "error", xErr)
		return xErr
	}

	if result.RowsAffected == 0 {
		r.Log.Warn(ctx, fmt.Sprintf("No rows affected when updating %s with id %v", r.table, id))
	}

	return nil
}
func (r *baseRepository[T]) UpdateWithIDs(ctx context.Context, ids []string, updatedData map[string]interface{}) error {
	op := "baseRepository.Update"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	// GORM: Update theo primary key
	updatedData["updated_at"] = gorm.Expr("CURRENT_TIMESTAMP") // Cập nhật thời gian cập nhật
	result := db.Table(r.table).Where("id in (?)", ids).Updates(updatedData)
	if result.Error != nil {
		xErr := apperrors.Wrap(result.Error, op, "Failed to update entity")
		r.Log.Error(ctx, "Failed to update entity", "error", xErr)
		return xErr
	}

	if result.RowsAffected == 0 {
		r.Log.Warn(ctx, fmt.Sprintf("No rows affected when updating %s with ids %v", r.table, ids))
	}

	return nil
}
func (r *baseRepository[T]) SoftDelete(ctx context.Context, id interface{}) error {
	op := "baseRepository.SoftDelete"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	var entity T
	if err := db.Where("id = ?", id).Delete(&entity).Error; err != nil {
		xErr := apperrors.Wrap(err, op, "Failed to soft delete entity")
		r.Log.Error(ctx, "Failed to soft delete entity", "error", xErr)
		return xErr
	}
	return nil
}
func (r *baseRepository[T]) HardDelete(ctx context.Context, id interface{}) error {
	op := "baseRepository.HardDelete"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	var entity T
	if err := db.Unscoped().Where("id = ?", id).Delete(&entity).Error; err != nil {
		xErr := apperrors.Wrap(err, op, "Failed to hard delete entity")
		r.Log.Error(ctx, "Failed to hard delete entity", "error", xErr)
		return xErr
	}
	return nil
}
func (r *baseRepository[T]) GetByID(ctx context.Context, id interface{}) (*T, error) {
	op := "baseRepository.GetByID"
	db := GetDb(ctx, r.db)
	db = db.WithContext(ctx)
	var entity T
	if err := db.Table(r.table).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		xErr := apperrors.Wrap(err, op, "Failed to get entity by ID")
		r.Log.Error(ctx, "Failed to get entity by ID", "error", xErr)
		return nil, xErr
	}
	return &entity, nil
}

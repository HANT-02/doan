package repositories

import (
	"context"

	apperrors "doan/pkg/error"
	"doan/pkg/logger"
)

// WorkFunc định nghĩa kiểu của hàm chứa logic nghiệp vụ
// sẽ được thực thi bên trong một transaction.
type WorkFunc func(txCtx context.Context) (interface{}, error)

// ExecuteInTransaction quản lý vòng đời transaction (Begin, Commit, Rollback)
// và phục hồi panic cho một unit of work.
func ExecuteInTransaction(
	ctx context.Context, // Context ban đầu, trước khi transaction bắt đầu
	uow UnitOfWork,
	log logger.Logger,
	fn WorkFunc,
) (result interface{}, err error) {
	op := "transaction.ExecuteInTransaction"
	var txCtx context.Context // Context dành riêng cho transaction

	txCtx, err = uow.Begin(ctx)
	if err != nil {
		err = apperrors.Wrap(err, op, "Failed to begin transaction")
		log.Error(ctx, "Failed to begin transaction", "error", err) // Log với context gốc
		return nil, err
	}

	committed := false
	defer func() {
		if r := recover(); r != nil {
			log.Error(txCtx, "Panic occurred during transaction", "error", r)
			if !committed {
				if rollbackErr := uow.Rollback(txCtx); rollbackErr != nil {
					log.Error(txCtx, "Failed to rollback transaction after panic", "error", rollbackErr)
				}
			}
			panic(r)
		} else if err != nil && !committed {
			log.Info(txCtx, "Transaction failed, rolling back", "triggering_error", err)
			if rollbackErr := uow.Rollback(txCtx); rollbackErr != nil {
				log.Error(txCtx, "Failed to rollback transaction after error", "error", rollbackErr)
			}
		}
	}()

	result, err = fn(txCtx)
	if err != nil {
		return nil, err
	}

	if err = uow.Commit(txCtx); err != nil {
		err = apperrors.Wrap(err, op, "Failed to commit transaction")
		log.Error(txCtx, "Failed to commit transaction", "error", err)
		return nil, err
	}
	committed = true

	return result, nil
}

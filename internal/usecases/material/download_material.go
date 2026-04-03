package material

import (
	"context"
	"errors"

	repositoryinterface "doan/internal/repositories/interface"
	auditservice "doan/internal/services/audit"
)

type DownloadMaterialInput struct {
	ID string
}

type DownloadMaterialOutput struct {
	ID           string
	FileName     string
	FileType     string
	FileSize     int64
	AbsolutePath string
}

type DownloadMaterialUseCase interface {
	Execute(ctx context.Context, input DownloadMaterialInput) (*DownloadMaterialOutput, error)
}

type downloadMaterialUseCase struct {
	materialRepo repositoryinterface.MaterialRepository
	storage      auditservice.StorageService
}

func NewDownloadMaterialUseCase(
	materialRepo repositoryinterface.MaterialRepository,
	storage auditservice.StorageService,
) DownloadMaterialUseCase {
	return &downloadMaterialUseCase{
		materialRepo: materialRepo,
		storage:      storage,
	}
}

func (uc *downloadMaterialUseCase) Execute(ctx context.Context, input DownloadMaterialInput) (*DownloadMaterialOutput, error) {
	if input.ID == "" {
		return nil, errors.New("material_id is required")
	}

	materialEntity, err := uc.materialRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if materialEntity == nil {
		return nil, errors.New("material not found")
	}

	resolvedFile, err := uc.storage.Resolve(ctx, materialEntity.FilePath)
	if err != nil {
		return nil, err
	}

	return &DownloadMaterialOutput{
		ID:           materialEntity.ID,
		FileName:     materialEntity.FileName,
		FileType:     materialEntity.FileType,
		FileSize:     materialEntity.FileSize,
		AbsolutePath: resolvedFile.AbsolutePath,
	}, nil
}

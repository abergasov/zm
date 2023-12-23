package filestree

import (
	"context"
	"fmt"
	"zm/internal/entities"
	"zm/internal/utils"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
)

func (r *Repository) SaveFiles(ctx context.Context, treeRoot string, files []*entities.FileMetadata) error {
	return r.saveFiles(ctx, r.db.Client(), treeRoot, files)
}

func (r *Repository) saveFiles(ctx context.Context, tx sqlx.ExecerContext, treeRoot string, files []*entities.FileMetadata) error {
	sql, params := utils.GenerateBulkInsertSQL(TableFiles, utils.PQParamPlaceholder, files, func(entity *entities.FileMetadata) map[string]any {
		return map[string]any{
			"f_id":       uuid.New(),
			"file_index": entity.FileIndex,
			"tree_id":    treeRoot,
			"file_hash":  entity.Hash,
			"file_name":  entity.FileName,
		}
	})
	_, err := tx.ExecContext(ctx, sql, params...)
	return err
}

func (r *Repository) GetFile(ctx context.Context, treeRoot string, fileIndex int) (*entities.FileMetadata, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE tree_id = $1 AND file_index = $2;", tableFilesFields, TableFiles)
	row := r.db.Client().QueryRowxContext(ctx, sql, treeRoot, fileIndex)
	var meta entities.FileMetadata
	err := row.StructScan(&meta)
	return &meta, err
}

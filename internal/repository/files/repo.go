package files

import (
	"context"
	"fmt"
	"strings"
	"zm/internal/entities"
	"zm/internal/storage/database"
	"zm/internal/utils"

	"github.com/google/uuid"
)

var (
	TableFiles       = "files"
	tableFilesFields = strings.Join([]string{"f_id", "file_index", "tree_id", "file_hash"}, ",")
)

type Repository struct {
	db database.DBConnector
}

func InitRepo(db database.DBConnector) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveFiles(ctx context.Context, treeID uuid.UUID, files []*entities.File) error {
	sql, params := utils.GenerateBulkInsertSQL(TableFiles, utils.PQParamPlaceholder, files, func(entity *entities.File) map[string]any {
		return map[string]any{
			"f_id":       uuid.New(),
			"file_index": entity.Meta.FileIndex,
			"tree_id":    treeID,
			"file_hash":  entity.Meta.Hash,
		}
	})
	_, err := r.db.Client().ExecContext(ctx, sql, params...)
	return err
}

func (r *Repository) GetFile(ctx context.Context, treeID uuid.UUID, fileIndex int) (*entities.FileMetadata, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE tree_id = $1 AND file_index = $2;", tableFilesFields, TableFiles)
	row := r.db.Client().QueryRowxContext(ctx, sql, treeID, fileIndex)
	var meta entities.FileMetadata
	err := row.StructScan(&meta)
	return &meta, err
}

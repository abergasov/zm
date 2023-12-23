package filestree

import (
	"context"
	"fmt"
	"strings"
	"zm/internal/entities"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/storage/database"
)

var (
	TableFiles       = "files"
	tableFilesFields = strings.Join([]string{
		"f_id",
		"file_index",
		"tree_id",
		"file_hash",
		"file_name",
	}, ",")

	TableTrees       = "merkle_trees"
	tableTreesFields = strings.Join([]string{"mt_id", "tree"}, ",") // nolint: gocritic
)

type Repository struct {
	db database.DBConnector
}

func InitRepo(db database.DBConnector) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveTreeAndFiles(ctx context.Context, tree *merkletree.Tree, files []*entities.FileMetadata) error {
	tx, err := r.db.Client().Beginx()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback() // nolint: errcheck
	if err = r.saveTree(ctx, tx, tree); err != nil {
		return fmt.Errorf("unable to save tree: %w", err)
	}
	if err = r.saveFiles(ctx, tx, tree.GetRoot(), files); err != nil {
		return fmt.Errorf("unable to save files: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}
	return nil
}

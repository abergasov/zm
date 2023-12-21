package tree

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/storage/database"

	"github.com/google/uuid"
)

var (
	TableTrees       = "merkle_trees"
	tableFilesFields = strings.Join([]string{"mt_id", "tree"}, ",")
)

type Repository struct {
	db database.DBConnector
}

func InitRepo(db database.DBConnector) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveTree(ctx context.Context, tree *merkletree.Tree) (uuid.UUID, error) {
	treeID := uuid.New()
	data, err := json.Marshal(tree)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to marshal tree: %w", err)
	}
	_, err = r.db.Client().ExecContext(ctx,
		fmt.Sprintf(`INSERT INTO %s (%s) VALUES ($1, $2);`, TableTrees, tableFilesFields), treeID, data,
	)
	return treeID, err
}

func (r *Repository) GetTree(ctx context.Context, treeID uuid.UUID) (*merkletree.Tree, error) {
	var treeBytes []byte
	err := r.db.Client().QueryRowContext(ctx, fmt.Sprintf(`SELECT tree FROM %s WHERE mt_id = $1;`, TableTrees), treeID).Scan(&treeBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to get tree: %w", err)
	}
	var tree *merkletree.Tree
	if err = json.Unmarshal(treeBytes, &tree); err != nil {
		return nil, fmt.Errorf("unable to unmarshal tree: %w", err)
	}
	return tree, nil
}

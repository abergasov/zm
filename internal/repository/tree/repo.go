package tree

import (
	"context"
	"encoding/json"
	"fmt"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/storage/database"
)

var (
	TableTrees          = "merkle_trees"
	tableClustersFields = []string{"mt_id", "tree"}
)

type Repository struct {
	db database.DBConnector
}

func InitRepo(db database.DBConnector) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveTree(ctx context.Context, tree *merkletree.Tree) (int64, error) {
	data, err := json.Marshal(tree)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal tree: %w", err)
	}
	res := r.db.Client().QueryRowContext(ctx, fmt.Sprintf(`INSERT INTO %s (tree) VALUES ($1) RETURNING mt_id;`, TableTrees), data)
	var treeID int64
	if err = res.Scan(&treeID); err != nil {
		return 0, fmt.Errorf("unable to scan tree id: %w", err)
	}
	return treeID, nil
}

func (r *Repository) GetTree(ctx context.Context, treeID int64) (*merkletree.Tree, error) {
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

package filestree

import (
	"context"
	"encoding/json"
	"fmt"
	merkletree "zm/internal/service/merkle_tree"

	"github.com/jmoiron/sqlx"
)

func (r *Repository) SaveTree(ctx context.Context, tree *merkletree.Tree) error {
	return r.saveTree(ctx, r.db.Client(), tree)
}

func (r *Repository) saveTree(ctx context.Context, tx sqlx.ExecerContext, tree *merkletree.Tree) error {
	treeID := tree.GetRoot()
	data, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("unable to marshal tree: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf(`INSERT INTO %s (%s) VALUES ($1, $2);`, TableTrees, tableTreesFields), treeID, data,
	)
	return err
}

func (r *Repository) GetTree(ctx context.Context, treeRoot string) (*merkletree.Tree, error) {
	var treeBytes []byte
	err := r.db.Client().QueryRowContext(ctx, fmt.Sprintf(`SELECT tree FROM %s WHERE mt_id = $1;`, TableTrees), treeRoot).Scan(&treeBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to get tree: %w", err)
	}
	var tree *merkletree.Tree
	if err = json.Unmarshal(treeBytes, &tree); err != nil {
		return nil, fmt.Errorf("unable to unmarshal tree: %w", err)
	}
	return tree, nil
}

package tree_test

import (
	"testing"
	merkletree "zm/internal/service/merkle_tree"
	testhelpers "zm/internal/test_helpers"
	"zm/internal/utils"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestRepo_CRUD(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	tree := merkletree.NewTree(utils.Hash256, uuid.NewString(), uuid.NewString(), uuid.NewString())

	// when
	require.NoError(t, container.RepositoryTrees.SaveTree(container.Ctx, tree))

	// then
	newTree, err := container.RepositoryTrees.GetTree(container.Ctx, tree.GetRoot())
	require.NoError(t, err)
	require.Equal(t, tree, newTree)
}

package filestree_test

import (
	"testing"
	"zm/internal/entities"
	merkletree "zm/internal/service/merkle_tree"
	testhelpers "zm/internal/test_helpers"
	"zm/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepoTrees_CRUD(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	tree := merkletree.NewTree(utils.Hash256, uuid.NewString(), uuid.NewString(), uuid.NewString())

	// when
	require.NoError(t, container.Repository.SaveTree(container.Ctx, tree))

	// then
	newTree, err := container.Repository.GetTree(container.Ctx, tree.GetRoot())
	require.NoError(t, err)
	require.Equal(t, tree, newTree)
}

func TestFiles_CRUD(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	tree := merkletree.NewTree(utils.Hash256, uuid.NewString(), uuid.NewString(), uuid.NewString())
	require.NoError(t, container.Repository.SaveTree(container.Ctx, tree))

	files := []*entities.FileMetadata{
		{
			FileIndex: 1,
			Hash:      uuid.NewString(),
			TreeRoot:  tree.GetRoot(),
		},
		{
			FileIndex: 2,
			Hash:      uuid.NewString(),
			TreeRoot:  tree.GetRoot(),
		},
		{
			FileIndex: 3,
			Hash:      uuid.NewString(),
			TreeRoot:  tree.GetRoot(),
		},
	}

	// when
	require.NoError(t, container.Repository.SaveFiles(container.Ctx, tree.GetRoot(), files))

	// then
	for _, file := range files {
		meta, err := container.Repository.GetFile(container.Ctx, tree.GetRoot(), file.FileIndex)
		require.NoError(t, err)
		file.FileID = meta.FileID
		require.Equal(t, file, meta)
	}
}

func TestRepository_SaveTreeAndFiles(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	tree := merkletree.NewTree(utils.Hash256, uuid.NewString(), uuid.NewString(), uuid.NewString())
	files := []*entities.FileMetadata{
		{
			FileIndex: 1,
			Hash:      uuid.NewString(),
			TreeRoot:  tree.GetRoot(),
		},
		{
			FileIndex: 2,
			Hash:      uuid.NewString(),
			TreeRoot:  tree.GetRoot(),
		},
		{
			FileIndex: 3,
			Hash:      uuid.NewString(),
			TreeRoot:  tree.GetRoot(),
		},
	}

	// when
	require.NoError(t, container.Repository.SaveTreeAndFiles(container.Ctx, tree, files))

	// then
	newTree, err := container.Repository.GetTree(container.Ctx, tree.GetRoot())
	require.NoError(t, err)
	require.Equal(t, tree, newTree)

	for _, file := range files {
		meta, err := container.Repository.GetFile(container.Ctx, tree.GetRoot(), file.FileIndex)
		require.NoError(t, err)
		file.FileID = meta.FileID
		require.Equal(t, file, meta)
	}
}

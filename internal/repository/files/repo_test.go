package files_test

import (
	"testing"
	"zm/internal/entities"
	merkletree "zm/internal/service/merkle_tree"
	testhelpers "zm/internal/test_helpers"
	"zm/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFiles_CRUD(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	tree := merkletree.NewTree(utils.Hash256, uuid.NewString(), uuid.NewString(), uuid.NewString())
	treeID, err := container.RepositoryTrees.SaveTree(container.Ctx, tree)
	require.NoError(t, err)

	files := []*entities.File{
		{
			Data: generateRandomBytes(),
			Meta: &entities.FileMetadata{
				FileIndex: 1,
				Hash:      uuid.NewString(),
				TreeID:    treeID,
			},
		},
		{
			Data: generateRandomBytes(),
			Meta: &entities.FileMetadata{
				FileIndex: 2,
				Hash:      uuid.NewString(),
				TreeID:    treeID,
			},
		},
		{
			Data: generateRandomBytes(),
			Meta: &entities.FileMetadata{
				FileIndex: 3,
				Hash:      uuid.NewString(),
				TreeID:    treeID,
			},
		},
	}

	// when
	require.NoError(t, container.RepositoryFiles.SaveFiles(container.Ctx, treeID, files))

	// then
	for _, file := range files {
		meta, err := container.RepositoryFiles.GetFile(container.Ctx, treeID, file.Meta.FileIndex)
		require.NoError(t, err)
		file.Meta.FileID = meta.FileID
		require.Equal(t, file.Meta, meta)
	}
}

func generateRandomBytes() []byte {
	res := make([]byte, 100)
	for i := range res {
		res[i] = byte(i)
	}
	return res
}

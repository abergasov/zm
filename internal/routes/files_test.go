package routes_test

import (
	"fmt"
	"testing"
	"zm/internal/entities"
	merkletree "zm/internal/service/merkle_tree"
	testhelpers "zm/internal/test_helpers"
	"zm/internal/utils"

	"github.com/stretchr/testify/require"
)

const (
	urlUploadFiles = "/api/v1/upload"
)

func Test_ProcessFiles(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	srv := testhelpers.NewTestServer(t, container)
	baseFolder := t.TempDir()
	for i := 0; i < 20; i++ {
		testhelpers.GenerateRandomFile(t, baseFolder)
	}
	tree, files, err := merkletree.CalculateTreeForFolder(baseFolder)
	require.NoError(t, err)
	t.Logf("generate files to upload, folder: %s, root: %s", baseFolder, tree.GetRoot())
	// when
	res := srv.UploadFolder(t, urlUploadFiles, utils.StringsFromObjectSlice(files, func(meta merkletree.FileMeta) string {
		return meta.Path
	}), entities.UploadFilesMeta{
		Root: tree.GetRoot(),
		Files: utils.StringsFromObjectSlice(files, func(meta merkletree.FileMeta) string {
			return meta.Name
		}),
	})

	// then
	res.RequireOk(t)
	t.Run("should serve file", func(t *testing.T) {
		for i := range files {
			// when
			res = srv.Get(t, fmt.Sprintf("/api/v1/file/%s/%d", tree.GetRoot(), i))

			// then
			res.RequireOk(t)
			var file entities.FileResponse
			res.RequireUnmarshal(t, &file)
			require.True(t, file.Proof.Verify(tree.GetRoot()))
			t.Logf("file %02d is valid, root: %s", i, tree.GetRoot())
		}
	})
}

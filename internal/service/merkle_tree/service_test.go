package merkletree_test

import (
	"encoding/json"
	"testing"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type tCase struct {
	name      string
	items     []string
	root      string
	expectLen int
}

func TestNewTree(t *testing.T) {
	t.Run("should serve error on empty items", func(t *testing.T) {
		// when
		tree := merkletree.NewTree(utils.Hash256)

		// then
		require.Equal(t, 0, tree.Len())
	})
	table := []tCase{
		{
			name:      "1 known items",
			items:     []string{"a"},
			root:      "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
			expectLen: 1,
		},
		{
			name:      "2 known items",
			items:     []string{"a", "b"},
			root:      "62af5c3cb8da3e4f25061e829ebeea5c7513c54949115b1acc225930a90154da",
			expectLen: 2,
		},
		{
			name:      "3 known items",
			items:     []string{"a", "b", "c"},
			root:      "35172c364a0d06a3ddbd3869ff682dd0395fad299787cda9c74cea0a14d8dc41",
			expectLen: 3,
		},
		{
			name:      "4 known items",
			items:     []string{"a", "b", "c", "d"},
			root:      "58c89d709329eb37285837b042ab6ff72c7c8f74de0446b091b6a0131c102cfd",
			expectLen: 3,
		},
		{
			name:      "7 known items",
			items:     []string{"a", "b", "c", "d", "e", "f", "g"},
			root:      "2fe9def8b466e7a3024800f39782e795ce5a712602eadb3f8de8c19fe26e8406",
			expectLen: 4,
		},
		{
			name:      "big items list",
			items:     generateItems(32767),
			expectLen: 16,
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			// when
			tree := merkletree.NewTree(utils.Hash256, tc.items...)

			// then
			checkTree(t, tree, tc)

			t.Run("restored tree check", func(t *testing.T) {
				// given
				data, err := json.Marshal(tree)
				require.NoError(t, err)

				// when
				var newTree *merkletree.Tree
				require.NoError(t, json.Unmarshal(data, &newTree))

				// then
				checkTree(t, newTree, tc)
			})
		})
	}

	t.Run("generate in loop", func(t *testing.T) {
		items := make([]string, 0, 100_000)
		for i := 1; i < 1_000; i++ {
			items = append(items, uuid.NewString())
			tree := merkletree.NewTree(utils.Hash256, items...)
			require.True(t, tree.Len() > 0)
		}
	})
}

func checkTree(t *testing.T, tree *merkletree.Tree, tc tCase) {
	require.Equal(t, tc.expectLen, tree.Len())
	if tc.root != "" {
		require.Equal(t, tc.root, tree.GetRoot())
	}

	t.Run("get proof for item", func(t *testing.T) {
		// given
		treeRoot := tree.GetRoot()

		for _, item := range tc.items {
			// when
			proof, err := tree.GetProofForItem(utils.Hash256(item))
			require.NoError(t, err)

			// then
			require.True(t, proof.Verify(treeRoot))
		}
	})

	t.Run("unknown item", func(t *testing.T) {
		// when
		proof, err := tree.GetProofForItem(utils.Hash256(uuid.NewString()))

		// then
		require.ErrorIs(t, err, merkletree.ErrItemNotFound)
		require.Nil(t, proof)
	})
}

func generateItems(count int) []string {
	items := make([]string, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, uuid.NewString())
	}
	return items
}

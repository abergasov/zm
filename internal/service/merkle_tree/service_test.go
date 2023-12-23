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
	//t.Skip()
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
			root:      "da3811154d59c4267077ddd8bb768fa9b06399c486e1fc00485116b57c9872f5",
			expectLen: 2,
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
		// given
		items := make([]string, 0, 1_000)
		for i := 0; i < 1_000; i++ {
			items = append(items, uuid.NewString())
			tree := merkletree.NewTree(utils.Hash256, items...)
			require.True(t, tree.Len() > 0)
			for _, item := range items {
				// when
				itemHash := utils.Hash256(item)
				proof, err := tree.GetProofForItem(itemHash)
				require.NoError(t, err)

				// then
				require.Truef(t, proof.Verify(tree.GetRoot()), "proof is not valid for item %d", i)
			}
		}
	})
}

func checkTree(t *testing.T, tree *merkletree.Tree, tc tCase) {
	require.Equal(t, tc.expectLen, tree.Len())
	expectItemsLen := len(tc.items)
	if expectItemsLen%2 != 0 {
		expectItemsLen++
	}
	require.Equal(t, expectItemsLen, tree.GetItemsLen())
	if tc.root != "" {
		require.Equal(t, tc.root, tree.GetRoot())
	}

	t.Run("get proof for item", func(t *testing.T) {
		// given
		badHash := utils.Hash256(uuid.NewString())
		treeRoot := tree.GetRoot()

		for _, item := range tc.items {
			// given
			proof, err := tree.GetProofForItem(utils.Hash256(item))
			require.NoError(t, err)
			require.True(t, proof.Verify(treeRoot))

			// when
			proof.Items[0][0] = badHash

			// then
			require.False(t, proof.Verify(tree.GetRoot()))
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

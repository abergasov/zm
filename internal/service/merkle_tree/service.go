package merkletree

import (
	"math"
	"zm/internal/utils"
)

type Tree struct {
	// each slice is a level of Tree. at the end we will have only one slice with root
	Tree     [][]string     `json:"tree"`
	ItemsMap map[string]int `json:"items_map"` // track items index by hash in Tree's [0]
}

func NewTree[T any](hasher func(T) string, items ...T) *Tree {
	if len(items) == 0 {
		return &Tree{
			Tree:     [][]string{},
			ItemsMap: map[string]int{},
		}
	}
	if len(items) == 1 {
		itemHash := hasher(items[0])
		return &Tree{
			Tree: [][]string{
				{
					itemHash,
					"",
				},
				{
					utils.HashItems(itemHash),
				},
			},
			ItemsMap: map[string]int{
				hasher(items[0]): 0,
			},
		}
	}
	t := &Tree{
		Tree:     make([][]string, merkleTreeDepth(len(items))),
		ItemsMap: make(map[string]int, len(items)),
	}
	t.Tree[0] = make([]string, 0, len(items))

	// generate leafs from items
	for i := range items {
		itemHash := hasher(items[i])
		t.Tree[0] = append(t.Tree[0], itemHash)
		t.ItemsMap[itemHash] = i
	}
	if len(items)%2 != 0 {
		t.Tree[0] = append(t.Tree[0], "")
	}

	offset := 0
	for len(t.Tree[offset]) > 1 {
		t.Tree[offset+1] = make([]string, 0, len(t.Tree[offset])/2) // each level has half of previous level

		for i := 0; i < len(t.Tree[offset]); i += 2 {
			hashes := make([]string, 0, 2)
			hashes = append(hashes, t.Tree[offset][i])
			if i+1 < len(t.Tree[offset]) {
				hashes = append(hashes, t.Tree[offset][i+1])
			}
			t.Tree[offset+1] = append(t.Tree[offset+1], utils.HashItems(hashes...))
		}
		if len(t.Tree[offset+1])%2 != 0 && len(t.Tree[offset+1]) != 1 {
			t.Tree[offset+1] = append(t.Tree[offset+1], "")
		}
		offset++
	}
	return t
}

// merkleTreeDepth calculates the depth of the Merkle tree.
func merkleTreeDepth(items int) int {
	if items == 0 {
		return 0
	}
	// Calculate the number of nodes in a binary Tree
	nodes := 2*items - 1
	// Calculate the depth of the binary Tree
	depth := int(math.Ceil(math.Log2(float64(nodes))))
	if depth == 0 {
		depth = 1
	}
	return depth
}

// Len returns the number of levels in the Merkle tree.
func (t *Tree) Len() int {
	return len(t.Tree)
}

// GetRoot returns the root hash of the Merkle tree.
func (t *Tree) GetRoot() string {
	if len(t.Tree) == 0 {
		return ""
	}
	return t.Tree[t.Len()-1][0]
}

// GetItemsLen returns the number of items in the Merkle tree.
func (t *Tree) GetItemsLen() int {
	if len(t.Tree) == 0 {
		return 0
	}
	return len(t.Tree[0])
}

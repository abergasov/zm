package merkletree

import (
	"errors"
	"zm/internal/utils"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type TreeProof struct {
	items [][]string
}

// GetProofForItem returns proof for item with given hash
// return hash of neighbor and path to root
func (t *Tree) GetProofForItem(itemHash string) (*TreeProof, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	index, ok := t.ItemsMap[itemHash]
	if !ok {
		return nil, ErrItemNotFound
	}
	proof := make([][]string, 0, t.Len())
	proof = append(proof, make([]string, 2))

	neighbor := index + 1
	if neighbor%2 == 0 {
		neighbor = index - 1
	}

	if index > neighbor {
		index, neighbor = neighbor, index
	}
	proof[0][0] = t.Tree[0][index]

	if neighbor < len(t.Tree[0]) {
		proof[0][1] = t.Tree[0][neighbor]
	}

	offset := 1
	for {
		proof = append(proof, make([]string, 0, 2))
		if len(t.Tree[offset]) == 1 {
			proof[offset] = append(proof[offset], t.Tree[offset][0])
			break
		}
		index = index / 2 // left child
		if index%2 != 0 {
			index = index - 1
		}
		neighbor = index + 1
		if neighbor == len(t.Tree[offset]) {
			neighbor = index - 1
		}
		if index > neighbor {
			index, neighbor = neighbor, index
		}
		proof[offset] = append(proof[offset], t.Tree[offset][index], t.Tree[offset][neighbor])
		offset++
	}
	return &TreeProof{items: proof}, nil
}

func (p *TreeProof) Verify(rootHash string) bool {
	if len(p.items) < 2 {
		return false
	}
	for i := range p.items {
		if len(p.items[i]) == 1 {
			return p.items[i][0] == rootHash
		}
		itemsHash := utils.HashItems(p.items[i]...)
		valid := itemsHash == p.items[i+1][0]
		if len(p.items[i+1]) == 2 {
			valid = valid || itemsHash == p.items[i+1][1]
		}
		if !valid {
			return false
		}
	}
	return false
}

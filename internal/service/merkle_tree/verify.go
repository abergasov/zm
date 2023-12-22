package merkletree

import (
	"errors"
	"zm/internal/utils"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type TreeProof struct {
	Items [][]string `json:"items"`
}

// GetProofForItem returns proof for item with given hash
// return hash of neighbor and path to root
func (t *Tree) GetProofForItem(itemHash string) (*TreeProof, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.Tree) == 1 {
		return &TreeProof{Items: [][]string{
			{t.Tree[0][0]},
		}}, nil
	}
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
		index /= 2 // left child
		if index%2 != 0 {
			index--
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
	return &TreeProof{Items: proof}, nil
}

func (p *TreeProof) Verify(rootHash string) bool {
	if len(p.Items) < 2 {
		return false
	}
	for i := range p.Items {
		if len(p.Items[i]) == 1 {
			return p.Items[i][0] == rootHash
		}
		itemsHash := utils.HashItems(p.Items[i]...)
		valid := itemsHash == p.Items[i+1][0]
		if len(p.Items[i+1]) == 2 {
			valid = valid || itemsHash == p.Items[i+1][1]
		}
		if !valid {
			return false
		}
	}
	return false
}

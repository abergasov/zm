package filer

import (
	"context"
	"fmt"
	"os"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/utils"
)

func (s *Service) ServeFile(ctx context.Context, treeRoot string, fileID int) (data []byte, proof *merkletree.TreeProof, err error) {
	tree, err := s.repoTree.GetTree(ctx, treeRoot)
	if err != nil {
		return nil, nil, fmt.Errorf("failed get tree: %w", err)
	}
	fileMeta, err := s.repoFile.GetFile(ctx, treeRoot, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed get file: %w", err)
	}
	proof, err = tree.GetProofForItem(fileMeta.Hash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed get proof: %w", err)
	}
	if !proof.Verify(treeRoot) {
		proof.Verify(treeRoot)
		return nil, nil, fmt.Errorf("proof is not valid")
	}
	// todo fix file reading
	filePrefix := utils.GetFormatString(tree.GetItemsLen())
	filePath := fmt.Sprintf("%s/%s/"+filePrefix+"_%s", s.filesFolder, treeRoot, fileMeta.FileIndex, fileMeta.FileName)
	data, err = os.ReadFile(filePath)
	return data, proof, err
}

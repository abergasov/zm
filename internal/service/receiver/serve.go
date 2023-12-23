package receiver

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/utils"
)

func (s *Service) ServeFile(ctx context.Context, treeRoot string, fileID int) (data []byte, proof *merkletree.TreeProof, err error) {
	l := s.log.With(slog.String("root", treeRoot), slog.Int("fileID", fileID))
	tree, err := s.repoFilesTrees.GetTree(ctx, treeRoot)
	if err != nil {
		l.Error("failed get tree", err)
		return nil, nil, fmt.Errorf("failed get tree: %w", err)
	}
	fileMeta, err := s.repoFilesTrees.GetFile(ctx, treeRoot, fileID)
	if err != nil {
		l.Error("failed get file", err)
		return nil, nil, fmt.Errorf("failed get file: %w", err)
	}
	proof, err = tree.GetProofForItem(fileMeta.Hash)
	if err != nil {
		l.Error("failed get proof", err, slog.String("hash", fileMeta.Hash), slog.String("treeRoot", treeRoot), slog.Int("fileID", fileID), slog.Int("treeItemsLen", tree.GetItemsLen()))
		return nil, nil, fmt.Errorf("failed get proof: %w", err)
	}
	if !proof.Verify(treeRoot) {
		l.Error("proof is not valid", fmt.Errorf("proof is not valid"))
		return nil, nil, fmt.Errorf("proof is not valid")
	}
	// todo fix file reading
	filePrefix := utils.GetFormatString(tree.GetItemsLen())
	filePath := fmt.Sprintf("%s/%s/"+filePrefix+"_%s", s.filesFolder, treeRoot, fileMeta.FileIndex, fileMeta.FileName)
	data, err = os.ReadFile(filePath)
	return data, proof, err
}

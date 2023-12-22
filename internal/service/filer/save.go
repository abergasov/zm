package filer

import (
	"context"
	"fmt"
	"zm/internal/entities"
	merkletree "zm/internal/service/merkle_tree"
)

func (s *Service) SaveFiles(ctx context.Context, meta *entities.UploadFilesMeta) error {
	tree, calculatedFiles, err := merkletree.CalculateTreeForFolder(fmt.Sprint(s.filesFolder, "/", meta.Root))
	if err != nil {
		return fmt.Errorf("failed calculate tree: %w", err)
	}
	if tree.GetRoot() != meta.Root {
		return fmt.Errorf("root hash is not equal")
	}
	if err = s.repoTree.SaveTree(ctx, tree); err != nil {
		s.log.Error("failed save tree", err)
		return fmt.Errorf("failed save tree: %w", err)
	}

	files := make([]*entities.FileMetadata, 0, len(meta.Files))
	for i := range meta.Files {
		files = append(files, &entities.FileMetadata{
			FileIndex: i,
			Hash:      calculatedFiles[i].Hash,
			FileName:  meta.Files[i],
			TreeRoot:  tree.GetRoot(),
		})
	}

	if err = s.repoFile.SaveFiles(ctx, tree.GetRoot(), files); err != nil {
		s.log.Error("failed save files", err)
		return fmt.Errorf("failed save files: %w", err)
	}
	return nil
}

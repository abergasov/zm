package receiver

import (
	"context"
	"fmt"
	"log/slog"
	"zm/internal/entities"
	merkletree "zm/internal/service/merkle_tree"
)

func (s *Service) SaveFiles(ctx context.Context, meta *entities.UploadFilesMeta) error {
	l := s.log.With(slog.String("root", meta.Root))
	l.Info("start save files")

	tree, calculatedFiles, err := merkletree.CalculateTreeForFolder(fmt.Sprint(s.filesFolder, "/", meta.Root))
	if err != nil {
		l.Error("failed calculate tree", err)
		return fmt.Errorf("failed calculate tree: %w", err)
	}

	if tree.GetRoot() != meta.Root {
		l.Error("root hash is not equal", fmt.Errorf("got: %s, expected: %s", tree.GetRoot(), meta.Root))
		return fmt.Errorf("root hash is not equal")
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

	if err = s.repoFilesTrees.SaveTreeAndFiles(ctx, tree, files); err != nil {
		l.Error("failed save files and tree", err)
		return fmt.Errorf("failed save files: %w", err)
	}
	l.Info("files saved")
	return nil
}

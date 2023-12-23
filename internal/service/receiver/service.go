package receiver

import (
	"os"
	"zm/internal/logger"
	filestree "zm/internal/repository/files_tree"
)

type Service struct {
	log            logger.AppLogger
	filesFolder    string
	repoFilesTrees *filestree.Repository
}

func NewReceiverService(log logger.AppLogger, repoFilesTrees *filestree.Repository, filesFolder string) *Service {
	if err := os.MkdirAll(filesFolder, os.ModePerm); err != nil {
		log.Fatal("unable to create storage folder", err)
	}
	return &Service{
		log:            log,
		filesFolder:    filesFolder,
		repoFilesTrees: repoFilesTrees,
	}
}

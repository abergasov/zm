package filer

import (
	"zm/internal/logger"
	"zm/internal/repository/files"
	"zm/internal/repository/tree"
)

type Service struct {
	log         logger.AppLogger
	filesFolder string
	repoTree    *tree.Repository
	repoFile    *files.Repository
}

func NewFilerService(log logger.AppLogger, repoTree *tree.Repository, repoFile *files.Repository, filesFolder string) *Service {
	return &Service{
		log:         log,
		filesFolder: filesFolder,
		repoTree:    repoTree,
		repoFile:    repoFile,
	}
}

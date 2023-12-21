package sampler

import (
	"log/slog"
	"zm/internal/logger"
	"zm/internal/repository/tree"
)

type Service struct {
	log  logger.AppLogger
	repo *tree.Repository
}

func InitService(log logger.AppLogger, repo *tree.Repository) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("service", "sampler")),
	}
}

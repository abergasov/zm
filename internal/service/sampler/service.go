package sampler

import (
	"log/slog"
	"zm/internal/logger"
	"zm/internal/repository/sampler"
)

type Service struct {
	log  logger.AppLogger
	repo *sampler.Repo
}

func InitService(log logger.AppLogger, repo *sampler.Repo) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("service", "sampler")),
	}
}

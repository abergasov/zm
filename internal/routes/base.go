package routes

import (
	"log/slog"
	"strings"
	"zm/internal/logger"
	"zm/internal/service/filer"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	appAddr      string
	filesFolder  string
	log          logger.AppLogger
	serviceFiles *filer.Service
	httpEngine   *fiber.App
}

// InitAppRouter initializes the HTTP Server.
func InitAppRouter(log logger.AppLogger, service *filer.Service, filesFolder, address string) *Server {
	app := &Server{
		appAddr:      address,
		filesFolder:  strings.TrimRight(filesFolder, "/"),
		httpEngine:   fiber.New(fiber.Config{}),
		serviceFiles: service,
		log:          log.With(slog.String("serviceFiles", "http")),
	}
	app.httpEngine.Use(recover.New())
	app.initRoutes()
	return app
}

func (s *Server) initRoutes() {
	s.httpEngine.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})
	s.httpEngine.Post("/api/v1/upload", s.handleFilesUpload)
	s.httpEngine.Get("/api/v1/file/:treeRoot/:fileID", s.serveFiles)
}

// Run starts the HTTP Server.
func (s *Server) Run() error {
	s.log.Info("Starting HTTP server", slog.String("port", s.appAddr))
	return s.httpEngine.Listen(s.appAddr)
}

func (s *Server) Stop() error {
	return s.httpEngine.Shutdown()
}

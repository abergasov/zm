package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"zm/internal/entities"
	"zm/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleFilesUpload(ctx *fiber.Ctx) error {
	metaRaw := ctx.FormValue("meta")
	if metaRaw == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	var meta entities.UploadFilesMeta
	if err := json.Unmarshal([]byte(metaRaw), &meta); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if len(meta.Files) == 0 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	l := s.log.With(slog.String("root", meta.Root))

	if err := os.MkdirAll(fmt.Sprintf("%s/%s", s.filesFolder, meta.Root), os.ModePerm); err != nil {
		l.Error("failed create folder for files", err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	filePrefix := utils.GetFormatString(len(meta.Files))
	for i := range meta.Files {
		file, err := ctx.FormFile(meta.Files[i])
		if err != nil {
			l.Error("failed fetch file from request", err, slog.String("file", meta.Files[i]))
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err = ctx.SaveFile(file, fmt.Sprintf("%s/%s/"+filePrefix+"_%s", s.filesFolder, meta.Root, i, file.Filename)); err != nil {
			l.Error("failed save file to tmp folder", err, slog.String("file", meta.Files[i]))
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}
	return s.serviceFiles.SaveFiles(ctx.Context(), &meta)
}

func (s *Server) serveFiles(ctx *fiber.Ctx) error {
	treeRoot := ctx.Params("treeRoot")
	if treeRoot == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	fileIDParams := ctx.Params("fileID")
	if fileIDParams == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	fileID, err := strconv.Atoi(fileIDParams)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	data, proof, err := s.serviceFiles.ServeFile(ctx.Context(), treeRoot, fileID)
	if err != nil {
		s.log.Error("failed serve file", err, slog.String("treeRoot", treeRoot), slog.Int("fileID", fileID))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(entities.FileResponse{
		Data:  data,
		Proof: proof,
	})
}

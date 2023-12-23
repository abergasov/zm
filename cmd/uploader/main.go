package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"zm/internal/entities"
	"zm/internal/logger"
	merkletree "zm/internal/service/merkle_tree"
	"zm/internal/utils"
)

var (
	confFile      = flag.String("config", "configs/app_conf.yml", "Configs file path")
	dataFolder    = flag.String("path", "data_folder", "folder with data to upload")
	serverAddress = flag.String("server", "http://127.0.0.1:8000", "server address")
)

func main() {
	flag.Parse()
	appLog := logger.NewAppSLogger("")
	appLog.Info("app starting", slog.String("conf", *confFile), slog.String("path", *dataFolder))

	tree, files, err := merkletree.CalculateTreeForFolder(*dataFolder)
	if err != nil {
		appLog.Fatal("unable to calculate tree", err)
	}
	appLog.Info("tree created", slog.String("root", tree.GetRoot()))
	appLog.Info("uploading files", slog.Int("count", len(files)))

	req, err := utils.CreateMultipartRequest(
		fmt.Sprintf("%s/api/v1/upload", *serverAddress),
		utils.StringsFromObjectSlice(files, func(meta merkletree.FileMeta) string {
			return meta.Path
		}),
		entities.UploadFilesMeta{
			Root: tree.GetRoot(),
			Files: utils.StringsFromObjectSlice(files, func(meta merkletree.FileMeta) string {
				return meta.Name
			}),
		},
	)
	if err != nil {
		appLog.Fatal("unable to create request", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		appLog.Fatal("unable to make request", err)
	}
	if res.StatusCode != http.StatusOK {
		appLog.Fatal("bad status code", fmt.Errorf("status code: %d", res.StatusCode))
	}
	_ = res.Body.Close()

	appLog.Info("files uploaded, check verification")

	for i, file := range files {
		l := appLog.With(slog.Int("index", i), slog.String("file", file.Name))
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/file/%s/%d", *serverAddress, tree.GetRoot(), i))
		if err != nil {
			l.Fatal("unable to make request", err)
		}
		if resp.StatusCode != http.StatusOK {
			l.Fatal("bad status code", fmt.Errorf("status code: %d", resp.StatusCode))
		}
		var response entities.FileResponse
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			l.Fatal("unable to decode response", err)
		}
		if !response.Proof.Verify(tree.GetRoot()) {
			l.Fatal("file is invalid", fmt.Errorf("file is invalid"))
		}
		_ = resp.Body.Close()
		hash, err := utils.GetFileHash(file.Path)
		if err != nil {
			l.Fatal("unable to get file hash", err)
		}
		if hash != utils.HashData(response.Data) {
			l.Fatal("file is invalid", fmt.Errorf("file is invalid"))
		}
		l.Info("file is valid")
	}
}

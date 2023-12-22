package main

import (
	"flag"
	"fmt"
	"log/slog"
	"zm/internal/logger"
	merkletree "zm/internal/service/merkle_tree"
)

var (
	confFile   = flag.String("config", "configs/app_conf.yml", "Configs file path")
	dataFolder = flag.String("path", "data_folder", "folder with data to upload")
	appLog     logger.AppLogger
)

func main() {
	for i := 0; i < 100; i++ {
		println(fmt.Sprintf("%04d", i))
	}
	flag.Parse()
	appLog = logger.NewAppSLogger("")
	appLog.Info("app starting", slog.String("conf", *confFile), slog.String("path", *dataFolder))

	tree, files, err := merkletree.CalculateTreeForFolder(*dataFolder)
	if err != nil {
		appLog.Fatal("unable to calculate tree", err)
	}
	appLog.Info("tree created", slog.String("root", tree.GetRoot()))
	appLog.Info("uploading files", slog.Int("count", len(files)))
}

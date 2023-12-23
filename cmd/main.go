package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zm/internal/config"
	"zm/internal/logger"
	"zm/internal/repository/files"
	"zm/internal/repository/tree"
	"zm/internal/routes"
	"zm/internal/service/filer"
	"zm/internal/storage/database"
)

var (
	confFile = flag.String("config", "configs/app_conf.yml", "Configs file path")
	appHash  = os.Getenv("GIT_HASH")
)

func main() {
	flag.Parse()
	appLog := logger.NewAppSLogger(appHash)

	appLog.Info("app starting", slog.String("conf", *confFile))
	appConf, err := config.InitConf(*confFile)
	if err != nil {
		appLog.Fatal("unable to init config", err, slog.String("config", *confFile))
	}

	appLog.Info("create storage connections")
	dbConn, err := getDBConnect(appLog, &appConf.ConfigDB, appConf.MigratesFolder)
	if err != nil {
		appLog.Fatal("unable to connect to db", err, slog.String("host", appConf.ConfigDB.Address))
	}
	defer func() {
		if err = dbConn.Close(); err != nil {
			appLog.Fatal("unable to close db connection", err)
		}
	}()

	appLog.Info("init repositories")
	repoTrees := tree.InitRepo(dbConn)
	repoFiles := files.InitRepo(dbConn)

	appLog.Info("init services")
	service := filer.NewFilerService(appLog, repoTrees, repoFiles, "/tmp")

	appLog.Info("init http service")
	appHTTPServer := routes.InitAppRouter(appLog, service, "/tmp", fmt.Sprintf(":%d", appConf.AppPort))
	defer func() {
		if err = appHTTPServer.Stop(); err != nil {
			appLog.Fatal("unable to stop http service", err)
		}
	}()
	go func() {
		if err = appHTTPServer.Run(); err != nil {
			appLog.Fatal("unable to start http service", err)
		}
	}()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
}

func getDBConnect(log logger.AppLogger, cnf *config.DBConf, migratesFolder string) (*database.DBConnect, error) {
	for i := 0; i < 5; i++ {
		dbConnect, err := database.InitDBConnect(cnf, migratesFolder)
		if err == nil {
			return dbConnect, nil
		}
		log.Error("can't connect to db", err, slog.Int("attempt", i))
		time.Sleep(time.Duration(i) * time.Second * 5)
	}
	return nil, fmt.Errorf("can't connect to db")
}

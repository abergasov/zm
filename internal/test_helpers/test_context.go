package testhelpers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"zm/internal/config"
	"zm/internal/logger"
	filestree "zm/internal/repository/files_tree"
	"zm/internal/service/receiver"
	"zm/internal/storage/database"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	testStorageFolder = "/tmp/zm"
)

type TestContainer struct {
	Ctx    context.Context
	Logger logger.AppLogger

	// repositories
	Repository *filestree.Repository

	// services deps
	ServiceFiles *receiver.Service
}

func GetClean(t *testing.T) *TestContainer {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	conf := getTestConfig(t)
	prepareTestDB(ctx, t, &conf.ConfigDB)

	dbConnect, err := database.InitDBConnect(&conf.ConfigDB, guessMigrationDir(t))
	require.NoError(t, err)
	cleanupDB(t, dbConnect)
	t.Cleanup(func() {
		require.NoError(t, dbConnect.Client().Close())
	})

	appLog := logger.NewAppSLogger("test")
	// repo init
	repoFilesTree := filestree.InitRepo(dbConnect)

	// service init
	serviceFiler := receiver.NewReceiverService(appLog, repoFilesTree, testStorageFolder)
	t.Cleanup(func() {
		cancel()
	})
	return &TestContainer{
		Ctx:    ctx,
		Logger: appLog,

		Repository: repoFilesTree,

		ServiceFiles: serviceFiler,
	}
}

func prepareTestDB(ctx context.Context, t *testing.T, cnf *config.DBConf) {
	dbConnect, err := database.InitDBConnect(&config.DBConf{
		Address:        cnf.Address,
		Port:           cnf.Port,
		User:           cnf.User,
		Pass:           cnf.Pass,
		DBName:         "postgres",
		MaxConnections: cnf.MaxConnections,
	}, "")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, dbConnect.Client().Close())
	}()
	if _, err = dbConnect.Client().ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", cnf.DBName)); !isDatabaseExists(err) {
		require.NoError(t, err)
	}
}

func getTestConfig(t *testing.T) *config.AppConfig {
	return &config.AppConfig{
		AppPort: 0,
		ConfigDB: config.DBConf{
			Address:        "localhost",
			Port:           "5449",
			User:           "aHAjeK",
			Pass:           "AOifjwelmc8dw",
			DBName:         getTestDB(t, "sybill_test"),
			MaxConnections: 10,
		},
	}
}

func isDatabaseExists(err error) bool {
	return checkSQLError(err, "42P04")
}

func checkSQLError(err error, code string) bool {
	if err == nil {
		return false
	}
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return string(pqErr.Code) == code
}

func guessMigrationDir(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	res := strings.Split(dir, "/internal")
	return res[0] + "/migrations"
}

func cleanupDB(t *testing.T, connector database.DBConnector) {
	tables := []string{
		filestree.TableFiles,
		filestree.TableTrees,
	}
	tx, err := connector.Client().Beginx()
	require.NoError(t, err)
	for _, table := range tables {
		_, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table))
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit())
}

func getTestDB(t *testing.T, dbName string) string {
	h := sha256.New()
	h.Write([]byte(t.Name()))
	bs := h.Sum(nil)
	testNameHash := fmt.Sprintf("%x", bs)
	maxLen := 7

	newDBName := fmt.Sprintf("%s_%s", dbName, strings.ToLower(testNameHash[:maxLen]))
	t.Log("test db name:", newDBName)
	return newDBName
}

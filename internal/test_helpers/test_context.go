package testhelpers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"zm/internal/config"
	"zm/internal/logger"
	"zm/internal/repository/tree"
	samplerService "zm/internal/service/sampler"
	"zm/internal/storage/database"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type TestContainer struct {
	Ctx context.Context

	// repositories
	RepositoryTrees *tree.Repository

	// services deps
	ServiceSampler *samplerService.Service
}

func GetClean(t *testing.T) *TestContainer {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	conf := getTestConfig()
	prepareTestDB(ctx, t, &conf.ConfigDB)

	dbConnect, err := database.InitDBConnect(&conf.ConfigDB, guessMigrationDir(t))
	require.NoError(t, err)
	cleanupDB(t, dbConnect)
	t.Cleanup(func() {
		require.NoError(t, dbConnect.Client().Close())
	})

	appLog := logger.NewAppSLogger("test")
	// repo init
	repoTree := tree.InitRepo(dbConnect)

	// service init
	serviceSampler := samplerService.InitService(appLog, repoTree)
	t.Cleanup(func() {
		cancel()
	})
	return &TestContainer{
		Ctx:             ctx,
		RepositoryTrees: repoTree,
		ServiceSampler:  serviceSampler,
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

func getTestConfig() *config.AppConfig {
	return &config.AppConfig{
		AppPort: 0,
		ConfigDB: config.DBConf{
			Address:        "localhost",
			Port:           "5449",
			User:           "aHAjeK",
			Pass:           "AOifjwelmc8dw",
			DBName:         "sybill_test",
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
		tree.TableTrees,
	}
	for _, table := range tables {
		_, err := connector.Client().Exec(fmt.Sprintf("TRUNCATE %s CASCADE", table))
		require.NoError(t, err)
	}
}

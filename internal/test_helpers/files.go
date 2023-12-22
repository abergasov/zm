package testhelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func GenerateRandomFile(t *testing.T, folder string) {
	fileName := uuid.NewString()
	filePath := filepath.Join(folder, fileName)
	file, err := os.Create(filePath)
	require.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString(uuid.NewString())
	require.NoError(t, err)
}

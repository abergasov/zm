package merkletree

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type FileMeta struct {
	Name string
	Path string
	Hash string
}

func CalculateTreeForFolder(folderPath string) (*Tree, []FileMeta, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read dir: %w", err)
	}
	fileList := make([]FileMeta, 0, len(files))
	fileHashMap := make(map[string]string, len(files))
	for i := range files {
		if files[i].IsDir() {
			continue
		}
		fContainer := FileMeta{
			Name: files[i].Name(),
			Path: fmt.Sprintf("%s/%s", folderPath, files[i].Name()),
		}
		fileHash, errH := getFileHash(fContainer.Path)
		if errH != nil {
			return nil, nil, fmt.Errorf("unable to get file hash: %w", errH)
		}
		fileHashMap[files[i].Name()] = fileHash
		fContainer.Hash = fileHash
		fileList = append(fileList, fContainer)
	}

	return NewTree(func(f FileMeta) string {
		return fileHashMap[f.Name]
	}, fileList...), fileList, err
}

func getFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", fmt.Errorf("unable to copy file: %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

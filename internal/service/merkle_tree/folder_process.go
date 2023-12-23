package merkletree

import (
	"fmt"
	"os"
	"zm/internal/utils"
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
		fileHash, errH := utils.GetFileHash(fContainer.Path)
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

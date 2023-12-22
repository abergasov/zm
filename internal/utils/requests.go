package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func CreateMultipartRequest(url string, filePaths []string, jsonData any) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add files to the request
	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close() // nolint: gocritic // test helper

		fileName := filepath.Base(filePath)
		part, err := writer.CreateFormFile(fileName, filePath)
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json: %w", err)
	}

	if err = writer.WriteField("meta", string(jsonBytes)); err != nil {
		return nil, fmt.Errorf("unable to write json: %w", err)
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, nil
}

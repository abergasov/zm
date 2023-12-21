package entities

import "github.com/google/uuid"

type File struct {
	Data []byte        `json:"data"`
	Meta *FileMetadata `json:"meta"`
}

type FileMetadata struct {
	FileID    uuid.UUID `json:"file_id" db:"f_id"`
	Hash      string    `json:"hash" db:"file_hash"`
	FileIndex int       `json:"file_index" db:"file_index"`
	TreeID    uuid.UUID `json:"tree_id" db:"tree_id"`
}

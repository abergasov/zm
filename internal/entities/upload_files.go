package entities

import merkletree "zm/internal/service/merkle_tree"

type UploadFilesMeta struct {
	Root  string   `json:"root"`
	Files []string `json:"files"`
}

type FileResponse struct {
	Data  []byte                `json:"data"`
	Proof *merkletree.TreeProof `json:"proof"`
}

package entities

type UploadFilesMeta struct {
	Root  string   `json:"root"`
	Files []string `json:"files"`
}

package types

type PrepareFileOutput struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	Checksum string `json:"checksum"`
	FileSize uint64 `json:"file_size"`
	FileType string `json:"file_type"`
}

type FileTypeInfo struct {
	DetectedMime string `json:"detected_mime"`
	Extension    string `json:"extension"`
}

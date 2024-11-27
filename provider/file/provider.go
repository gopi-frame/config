package file

import (
	"os"
)

type FileProvider struct {
	filename string
}

func NewFileProvider(filename string) *FileProvider {
	return &FileProvider{
		filename: filename,
	}
}

func (fp *FileProvider) Read() ([]byte, error) {
	return os.ReadFile(fp.filename)
}

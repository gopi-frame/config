package fsys

import "io/fs"

type FSProvider struct {
	fsys     fs.FS
	filename string
}

func NewFSProvider(fsys fs.FS, filename string) *FSProvider {
	return &FSProvider{fsys: fsys, filename: filename}
}

func (p *FSProvider) Read() ([]byte, error) {
	return fs.ReadFile(p.fsys, p.filename)
}

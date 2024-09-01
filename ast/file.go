package ast

import (
	"io/fs"
	"os"
)

type File struct {
	Scope    *Scope
	FileName string
	FileInfo os.FileInfo
	fsys     fs.FS
}

// NewFile stat the file to get the file info
func NewFile(fsys fs.FS, filename string) (*File, error) {
	fi, err := fs.Stat(fsys, filename)
	if err != nil {
		return nil, err
	}
	return &File{FileName: filename, FileInfo: fi, fsys: fsys}, nil
}

func (f *File) ReadFile() ([]byte, error) {
	return fs.ReadFile(f.fsys, f.FileName)
}

func (f *File) String() string {
	return f.FileName
}

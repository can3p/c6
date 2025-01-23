package ast

import (
	"strings"
)

/*
For import like this:

@import "component/list"; // => component/_list.scss
*/
type ImportStmt struct {
	SourceFileName string
	Paths          []*String
}

func NewImportStmt(sourceFileName string) *ImportStmt {
	return &ImportStmt{
		SourceFileName: sourceFileName,
	}
}

func (self ImportStmt) CanBeStmt() {}

func (self ImportStmt) String() string {
	var b strings.Builder
	b.WriteString("@import ")
	// we expect at least one
	for idx, p := range self.Paths {
		if idx > 0 {
			b.WriteString(", ")
		}

		b.WriteString(p.String())
	}
	b.WriteString(";")

	return b.String()
}

package parser

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/c9s/c6/ast"
)

const (
	UnknownFileType uint = iota
	ScssFileType
	SassFileType
	EcssFileType
)

func debug(format string, args ...interface{}) {
	if debugParser {
		fmt.Printf(format+"\n", args...)
	}
}

func getFileTypeByExtension(extension string) uint {
	switch extension {
	case ".scss":
		return ScssFileType
	case ".sass":
		return SassFileType
	case ".ecss":
		return EcssFileType
	}
	return UnknownFileType
}

// one day we'll put some options there
// Parser represents a global instance of parser that
// contains shared options for all the files
type GlobalParser struct {
	fsys fs.FS
}

// Parser represent the state of parsing of a given file
type Parser struct {
	GlobalParser *GlobalParser

	File *ast.File

	// file content
	Content string

	// Integer for counting token
	Pos int

	// Position saved for rollbacking back.
	RollbackPos int

	// A token slice that contains all lexed tokens
	Tokens []*ast.Token
}

func NewParser(fsys fs.FS) *GlobalParser {
	return &GlobalParser{
		fsys: fsys,
	}
}

func (gp *GlobalParser) ResolveFileFname(source string, name string) (string, error) {
	base := path.Dir(source)
	importPath := path.Join(base, name)

	fi, err := fs.Stat(gp.fsys, importPath)

	if err != nil && os.IsExist(err) {
		return "", err
	}

	// @import "foo" may import foo.scss from the current folder
	if !strings.Contains(name, "/") {
		currentFolderFile := importPath + ".scss"

		fi, err := fs.Stat(gp.fsys, currentFolderFile)

		if err != nil && os.IsExist(err) {
			return "", err
		}

		if fi != nil {
			return currentFolderFile, nil
		}
	}

	// go find the _index.scss if it's a local directory
	if fi != nil && fi.Mode().IsDir() {
		importPath = path.Join(importPath, "_index.scss")
	} else {
		var dirname = filepath.Dir(importPath)
		var basename = filepath.Base(importPath)
		importPath = path.Join(dirname, "_"+basename+".scss")
	}

	return importPath, nil
}

func (gp *GlobalParser) ParseFile(path string) (*ast.StmtList, error) {
	ext := filepath.Ext(path)
	filetype := getFileTypeByExtension(ext)

	f, err := ast.NewFile(gp.fsys, path)
	if err != nil {
		return nil, err
	}
	data, err := f.ReadFile()
	if err != nil {
		return nil, err
	}

	parser := &Parser{
		GlobalParser: gp,
		Content:      string(data),
		File:         f,
	}

	var stmts *ast.StmtList

	switch filetype {
	case ScssFileType:
		stmts, err = parser.ParseScss(parser.Content)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported file format: %s", path)
	}
	return stmts, nil
}

func (gp *GlobalParser) ParseScss(content string) (*ast.StmtList, error) {
	parser := &Parser{
		GlobalParser: gp,
		Content:      content,
	}

	return parser.ParseScss(parser.Content)
}

func (parser *Parser) backup() {
	parser.Pos--
}

func (parser *Parser) restore(pos int) {
	parser.Pos = pos
}

//func (parser *Parser) remember() {
//parser.RollbackPos = parser.Pos
//}

// rollback to the save position
//func (parser *Parser) rollback() {
//parser.Pos = parser.RollbackPos
//}

// accept() accepts one token type one time.
// rolls back if the token type mismatch
func (parser *Parser) accept(tokenType ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok != nil && tok.Type == tokenType {
		return tok
	}
	parser.backup()
	return nil
}

// acceptAny accepts some token types, or it rolls back when the token mismatch
// the token types.
//func (parser *Parser) acceptAny(tokenTypes ...ast.TokenType) *ast.Token {
//var tok = parser.next()
//for _, tokType := range tokenTypes {
//if tok.Type == tokType {
//return tok
//}
//}
//parser.backup()
//return nil
//}

func (parser *Parser) acceptAnyOf2(tokType1, tokType2 ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok.Type == tokType1 || tok.Type == tokType2 {
		return tok
	}
	parser.backup()
	return nil
}

func (parser *Parser) acceptAnyOf3(tokType1, tokType2, tokType3 ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok.Type == tokType1 || tok.Type == tokType2 || tok.Type == tokType3 {
		return tok
	}
	parser.backup()
	return nil
}

func (parser *Parser) expect(tokenType ast.TokenType) (*ast.Token, error) {
	var tok = parser.next()
	if tok != nil && tok.Type != tokenType {
		parser.backup()
		return nil, SyntaxError{
			Reason:      fmt.Sprintf("expected: %s", tokenType.String()),
			ActualToken: tok,
			File:        parser.File,
		}
	}
	return tok, nil
}

func (parser *Parser) next() *ast.Token {
	var p = parser.Pos
	parser.Pos++

	// if we've appended the token
	if p < len(parser.Tokens) {
		return parser.Tokens[p]
	}
	return nil
}

func (parser *Parser) peekBy(offset int) *ast.Token {
	var i = 0
	var tok *ast.Token = nil
	for offset > 0 {
		tok = parser.next()
		offset--
		i++
		if tok == nil {
			break
		}
	}
	parser.Pos -= i
	return tok
}

func (parser *Parser) advance() {
	parser.Pos++
}

//func (parser *Parser) current() *ast.Token {
//return parser.Tokens[parser.Pos]
//}

func (parser *Parser) peek() *ast.Token {
	if parser.Pos < len(parser.Tokens) {
		return parser.Tokens[parser.Pos]
	}
	return nil
}

func (parser *Parser) eof() bool {
	var tok = parser.next()
	parser.backup()
	return tok == nil
}

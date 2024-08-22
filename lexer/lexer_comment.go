package lexer

import (
	"github.com/c9s/c6/ast"
)

func lexCommentLine(l *Lexer) (stateFn, error) {
	if !l.match("//") {
		return nil, nil
	}
	l.ignore()

	var r = l.next()
	for r != EOF {
		if r == '\n' {
			break
		}
		r = l.next()
		if r == '\r' {
			r = l.next()
		}
	}
	l.backup()
	l.ignore()
	return lexStart, nil
}

func lexCommentBlock(l *Lexer, emit bool) (stateFn, error) {
	if !l.match("/*") {
		return nil, nil
	}
	l.ignore()
	var r = l.next()
	for r != EOF {
		if r == '*' && l.peek() == '/' {
			l.backup()
			if emit {
				l.emit(ast.T_COMMENT_BLOCK)
			} else {
				l.ignore()
			}
			l.match("*/")
			l.ignore()
			return lexStart, nil
		}
		r = l.next()
	}
	return nil, l.errorf("Expecting comment end mark '*/'. Got '%c'", r)
}

func lexComment(l *Lexer, emit bool) (stateFn, error) {
	var r = l.peek()
	var r2 = l.peekBy(2)
	if r == '/' && r2 == '*' {
		return lexCommentBlock(l, emit)
	} else if r == '/' && r2 == '/' {
		return lexCommentLine(l)
	}
	return nil, nil
}

package lexer

import (
	"fmt"
	"unicode"

	"github.com/c9s/c6/ast"
)

func lexIdentifier(l *Lexer) (stateFn, error) {
	var r = l.next()
	if !unicode.IsLetter(r) && r != '-' {
		return nil, fmt.Errorf("An identifier needs to start with a letter or dash")
	}
	r = l.next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
		if r == '-' {
			var r2 = l.peek()
			if !unicode.IsLetter(r2) && r2 != '-' {
				l.backup()
				return lexExpr, nil
			}
		}

		r = l.next()
	}
	l.backup()

	if l.peek() == '(' {
		var curTok = l.emit(ast.T_FUNCTION_NAME)

		if curTok.Str == "url" || curTok.Str == "local" {
			if err := lexUrlParam(l); err != nil {
				return nil, err
			}
		} else {
			if _, err := lexFunctionParams(l); err != nil {
				return nil, err
			}
		}
	} else {
		l.emit(ast.T_IDENT)
	}
	return lexExpr, nil
}

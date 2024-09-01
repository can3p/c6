package lexer

import (
	"unicode"

	"github.com/c9s/c6/ast"
)

// $var-rgba(255,255,0)
func lexVariableName(l *Lexer) (stateFn, error) {
	var r = l.next()
	if r != '$' {
		return nil, l.errorf("Unexpected token %c for lexVariable", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		return nil, l.errorf("The first character of a variable name must be letter. Got '%c'", r)
	}

	r = l.next()
	for r != EOF {
		if r == '-' {
			var r2 = l.peek()
			if unicode.IsLetter(r2) { // $a-b is a valid variable name.
				l.next()
			} else if unicode.IsDigit(r2) { // $a-3 should be $a '-' 3
				l.backup()
				l.emit(ast.T_VARIABLE)
				return lexExpr, nil
			} else {
				break
			}
		} else if r == ':' {
			break
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		} else if r == '}' {
			l.backup()
			l.emit(ast.T_VARIABLE)
			return lexStart, nil
			///XXX break
		} else if unicode.IsSpace(r) || r == ';' {
			break
		}
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_VARIABLE)

	if l.match("...") {
		l.emit(ast.T_VARIABLE_LENGTH_ARGUMENTS)
	}

	return lexStart, nil
}

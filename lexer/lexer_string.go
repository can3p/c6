package lexer

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func lexString(l *Lexer) (stateFn, error) {
	var r = l.next()
	if r == '"' {
		var containsInterpolation = false
		l.ignore()
		// string start
		r = l.next()
		for {
			if r == '"' {
				l.backup()
				token := l.createToken(ast.T_QQ_STRING)
				token.ContainsInterpolation = containsInterpolation
				l.emitToken(token)
				l.next()
				l.ignore()
				return lexStart, nil
			} else if r == '\\' {
				// skip the escape character
			} else if IsInterpolationStartToken(r, l.peek()) {
				l.backup()
				if _, err := lexInterpolation(l, false); err != nil {
					return nil, err
				}
				containsInterpolation = true
			} else if r == EOF {
				return nil, fmt.Errorf("Expecting end of string")
			}
			r = l.next()
		}
		//XXX l.backup()
		//XXX return lexStart

	} else if r == '\'' {
		l.ignore()
		l.next()
		for {
			r = l.next()
			if r == '\'' {
				l.backup()
				l.emit(ast.T_Q_STRING)
				l.next()
				l.ignore()
				return lexStart, nil
			} else if r == '\\' {
				// skip the escape character
				l.next()
			} else if r == EOF {
				return nil, fmt.Errorf("Expecting end of string")
			}
		}
		//XXX return lexStart
	}
	l.backup()
	return nil, nil
}

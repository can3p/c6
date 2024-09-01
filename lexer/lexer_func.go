package lexer

import (
	"github.com/c9s/c6/ast"
)

func lexFunctionParams(l *Lexer) (stateFn, error) {
	var r = l.next()
	if r != '(' {
		return nil, l.errorf("Expecting '('. Got '%c'.", r)
	}
	l.emit(ast.T_PAREN_OPEN)
	l.ignoreSpaces()

	r = l.peek()
	for r != EOF {
		l.ignoreSpaces()
		r = l.peek()
		if r == ')' {
			l.next()
			l.emit(ast.T_PAREN_CLOSE)
			break
		}
		if expr, err := lexExpr(l); err != nil {
			return nil, err
		} else if expr == nil {
			break
		}

		l.ignoreSpaces()
		r = l.peek()
		if r == ',' {
			l.next()
			l.emit(ast.T_COMMA)
		}
	}
	return nil, nil
}

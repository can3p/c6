package lexer

import (
	"unicode"

	"github.com/c9s/c6/ast"
)

/*
There are 3 scope that users may use interpolation syntax:

	 {selector interpolation}  {
		 {property name inpterpolation}: {property value interpolation}
	 }
*/
func lexInterpolation(l *Lexer, emit bool) (stateFn, error) {
	l.remember()
	var r rune = l.next()
	if r == '#' {
		r = l.next()
		if r == '{' {
			if emit {
				l.emit(ast.T_INTERPOLATION_START)
			}

			r = l.next()
			for unicode.IsSpace(r) {
				r = l.next()
			}
			l.backup()

			// find the end of interpolation end brace
			r = l.next()
			for r != '}' {
				r = l.next()
			}
			l.backup()

			if emit {
				l.emit(ast.T_INTERPOLATION_INNER)
			}

			l.next() // for '}'
			if emit {
				l.emit(ast.T_INTERPOLATION_END)
			}
			return nil, nil
		}
	}
	l.rollback()
	return nil, nil
}

// Lex the expression inside interpolation
func lexInterpolation2(l *Lexer) (stateFn, error) {
	var r rune = l.next()
	if r != '#' {
		return nil, l.errorf("Expecting interpolation token '#', Got %c", r)
	}
	r = l.next()
	if r != '{' {
		return nil, l.errorf("Expecting interpolation token '{', Got %c", r)
	}
	l.emit(ast.T_INTERPOLATION_START)

	// skip the space after #{
	for {
		expr, err := lexExpr(l)
		if err != nil {
			return nil, err
		}

		if expr == nil {
			break
		}
	}
	l.ignoreSpaces()
	if err := l.expect("}"); err != nil {
		return nil, err
	}
	l.emit(ast.T_INTERPOLATION_END)
	return nil, nil
}

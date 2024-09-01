package lexer

import (
	"unicode"

	"github.com/c9s/c6/ast"
)

func lexPropertyNameToken(l *Lexer) (stateFn, error) {
	var r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
		r = l.next()
	}
	l.backup()
	if l.precedeStartOffset() {
		l.emit(ast.T_PROPERTY_NAME_TOKEN)
		return lexPropertyNameToken, nil
	}
	return nil, nil
}

func lexMicrosoftProgIdFunction(l *Lexer) (stateFn, error) {
	var r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_FUNCTION_NAME)

	// here starts the sproperty
	r = l.next()
	if r != '(' {
		return nil, l.errorf("Expecting '(' after the MS function name. Got %c", r)
	}
	l.emit(ast.T_PAREN_OPEN)

	l.ignoreSpaces()
	if err := l.ignoreComment(); err != nil {
		return nil, err
	}

	r = l.peek()
	l.ignoreSpaces()

	// here comes the sProperty
	//     progid:DXImageTransform.Microsoft.filtername(sProperties)
	// @see https://msdn.microsoft.com/en-us/library/ms532847(v=vs.85).aspx
	for r != ')' {
		// lex function parameter name
		if unicode.IsSpace(r) {
			l.ignoreSpaces()
		}

		if err := l.ignoreComment(); err != nil {
			return nil, err
		}

		r = l.next()
		for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			r = l.next()
		}
		l.backup()
		l.emit(ast.T_MS_PARAM_NAME)
		if l.accept("=") {
			l.emit(ast.T_ATTR_EQUAL)

			if _, err := lexExpr(l); err != nil {
				return nil, err
			}

			l.ignoreSpaces()
			r = l.peek()
			if r == ',' {
				l.next()
				l.emit(ast.T_COMMA)
				l.ignoreSpaces()
			} else if r == ')' {
				l.next()
				l.emit(ast.T_PAREN_CLOSE)
				break
			}
		} else if r == ')' {
			l.next()
			l.emit(ast.T_PAREN_CLOSE)
			break
		}
	}
	return nil, nil
}

/*
Possible property value syntax:

	width: 10px;        // numeric
	width: 10px + 10px; // expression
	border: 1px #{solid} #000;   // interpolation
	color: rgba( 0, 0, 255, 1.0);  // rgba function
	width: auto;    // string constant
*/
func lexProperty(l *Lexer) (stateFn, error) {
	var r = l.peek()

	for r != ':' && r != '/' && !unicode.IsSpace(r) {
		if r == '.' {
			return nil, l.errorf("dot notation in properties is not supported. Got %c", r)
		}

		if l.peek() == '#' && l.peekBy(2) == '{' {
			if _, err := lexInterpolation2(l); err != nil {
				return nil, err
			}

			r = l.peek()
			if !unicode.IsSpace(r) && r != ':' {
				l.emit(ast.T_LITERAL_CONCAT)
			}
		}

		// we have something
		if t, err := lexPropertyNameToken(l); err != nil {
			return nil, err
		} else if t != nil {
			r = l.peek()
			if !unicode.IsSpace(r) && r != ':' {
				l.emit(ast.T_LITERAL_CONCAT)
			}
		}
		r = l.peek()
	}

	l.ignoreSpaces()
	if l.peek() == '/' {
		if _, err := lexComment(l, false); err != nil {
			return nil, err
		}
		l.ignoreSpaces()
	}

	if _, err := lexColon(l); err != nil {
		return nil, err
	}

	l.remember()
	l.ignoreSpaces()

	// for IE filter syntax like:
	//    progid:DXImageTransform.Microsoft.MotionBlur(strength=13, direction=310)
	if l.match("progid:") {
		l.emit(ast.T_MS_PROGID)
		_, err := lexMicrosoftProgIdFunction(l)
		if err != nil {
			return nil, err
		}
	} else {
		l.rollback()
	}

	// the '{' is used for the start token of nested properties
	r = l.peek()
	for r != EOF {

		// for nested property value
		if r == '{' {

			l.next()
			l.emit(ast.T_BRACE_OPEN)
			return lexStart, nil

			//} else if lexExpr(l) == nil {
		} else if expr, err := lexExpr(l); err != nil {
			return nil, err
		} else if expr == nil {
			break
		}

		r = l.peek()
	}

	l.ignoreSpaces()
	if _, err := lexComment(l, false); err != nil {
		return nil, err
	}
	l.ignoreSpaces()

	// the semicolon in the last declaration is optional.
	l.ignoreSpaces()
	if l.accept(";") {
		l.emit(ast.T_SEMICOLON)
	}

	l.ignoreSpaces()
	if l.accept("}") {
		l.emit(ast.T_BRACE_CLOSE)
	}
	return lexStart, nil
}

func lexColon(l *Lexer) (stateFn, error) {
	l.ignoreSpaces()
	var r = l.next()
	if r != ':' {
		return nil, l.errorf("Expecting ':' token, Got '%c'", r)
	}
	l.emit(ast.T_COLON)

	// We don't ignore space after the colon because we need spaces to detect literal concat.
	return nil, nil
}

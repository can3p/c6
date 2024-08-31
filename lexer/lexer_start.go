package lexer

import (
	"fmt"
	"unicode"

	"github.com/c9s/c6/ast"
)

func lexStart(l *Lexer) (stateFn, error) {
	// strip the leading spaces of a statement
	l.ignoreSpaces()

	var r, r2 = l.peek2()

	// lex simple statements
	switch r {

	case EOF:
		return nil, nil

	case '@':
		return lexAtRule, nil
	case '(':
		l.next()
		l.emit(ast.T_PAREN_OPEN)
		return lexStart, nil
	case ')':
		l.next()
		l.emit(ast.T_PAREN_CLOSE)
		return lexStart, nil

	case '{':
		l.next()
		l.emit(ast.T_BRACE_OPEN)
		return lexStart, nil

	case '}':
		l.next()
		l.emit(ast.T_BRACE_CLOSE)
		return lexStart, nil

	case '$':
		return lexAssignStmt, nil

	case ';':
		l.next()
		l.emit(ast.T_SEMICOLON)
		return lexStart, nil

	case '-':
		// css/custom_properties/indentation.hrx
		// lexProperty gets into endless loop while trying to parse
		// a css variable.
		// @TODO: remove and immplement css variables support properly
		if r2 == '-' {
			return nil, fmt.Errorf("CSS variables are not supported yet")
		}
		// Vendor prefix properties start with '-'
		return lexProperty, nil

	case ',':
		l.next()
		l.emit(ast.T_COMMA)
		return lexStart, nil

	case '/':
		if r2 == '*' {

			if _, err := lexCommentBlock(l, true); err != nil {
				return nil, err
			}

			return lexStart, nil

		} else if r2 == '/' {

			return lexCommentLine, nil

		}

		return nil, fmt.Errorf("unexpected token. expecing '*' or '/'")

	case '#':
		// make sure it's not an interpolation "#{" token
		if r2 != '{' {
			// looks like a selector
			return lexSelectors, nil
		}

	case '[', '*', '>', '&', '.', '+', ':':
		return lexSelectors, nil

	}

	if l.match("<!--") {

		l.emit(ast.T_CDOPEN)

		return lexStart, nil
	}

	if l.match("-->") {

		l.emit(ast.T_CDCLOSE)

		return lexStart, nil
	}

	// If a line starts with a letter or a sharp,
	// it might be a property name or selector, e.g.,
	//
	//    ul { }
	//
	//    -webkit-border-radius: 3px;
	//
	if unicode.IsLetter(r) || (r == '#') { // it might be -vendor- property or a property name or a selector

		// detect selector syntax
		l.remember()

		isProperty := false

		r = l.next()
		for r != EOF {
			// skip interpolation
			if r == '#' {
				if l.peek() == '{' {
					// find the matching brace
					r = l.next()
					for r != '}' {
						r = l.next()
					}
				}

			} else if r == ':' { // pseudo selector -> letters following ':', if there is a space after the ':' then it's a property value.

				if unicode.IsSpace(l.peek()) {
					isProperty = true
					break
				}

			} else if r == ';' {
				break
			} else if r == '}' {
				isProperty = true
				break
			} else if r == '{' {
				break
			} else if r == EOF {
				return nil, fmt.Errorf("unexpected EOF")
			}
			r = l.next()
		}

		l.rollback()

		if isProperty {
			return lexProperty, nil
		} else {
			return lexSelectors, nil
		}
	}
	return nil, l.errorf("Unexpected token: '%c'", r)
}

package lexer

import (
	"testing"

	"github.com/c9s/c6/ast"
	"github.com/stretchr/testify/assert"
)

func AssertLexerTokenSequenceFromState(t *testing.T, scss string, fn stateFn, tokenList []ast.TokenType) {
	t.Logf("Testing SCSS: %s\n", scss)
	var lexer = NewLexerWithString(scss)
	assert.NotNil(t, lexer)
	_, err := lexer.RunFrom(fn)
	assert.NoError(t, err)
	AssertTokenSequence(t, lexer, tokenList)
}

func AssertLexerTokenSequence(t *testing.T, scss string, tokenList []ast.TokenType) {
	t.Logf("Testing SCSS: %s\n", scss)
	var lexer = NewLexerWithString(scss)
	assert.NotNil(t, lexer)
	_, err := lexer.Run()
	assert.NoError(t, err)
	AssertTokenSequence(t, lexer, tokenList)
}

func OutputGreen(t *testing.T, msg string, args ...interface{}) {
	t.Logf("\033[32m")
	t.Logf(msg, args...)
	t.Logf("\033[0m\n")
}

func OutputRed(t *testing.T, msg string, args ...interface{}) {
	t.Logf("\033[31m")
	t.Logf(msg, args...)
	t.Logf("\033[0m\n")
}

func AssertTokenSequence(t *testing.T, l *Lexer, tokenList []ast.TokenType) []ast.Token {

	var tokens = []ast.Token{}
	var failure = false
	for idx, expectingToken := range tokenList {

		if idx == len(l.Tokens) {
			for _, token := range tokenList[idx:] {
				OutputRed(t, "not ok ---- input longer than expected  %s", token.String())
				failure = true
			}
			break
		}

		var token = l.Tokens[idx]

		if token == nil {
			failure = true
			OutputRed(t, "not ok ---- got nil expecting %s", expectingToken.String())
			break
		}

		tokens = append(tokens, *token)

		if expectingToken == token.Type {
			//OutputGreen(t, "ok %s '%s'", token.Type.String(), token.Str)
		} else {
			failure = true
			OutputRed(t, "not ok ---- %d token => got %s '%s' expecting %s", idx, token.Type.String(), token.Str, expectingToken.String())
		}
		assert.Equal(t, expectingToken, token.Type)
	}

	if len(tokenList) < len(l.Tokens) {
		for _, token := range l.Tokens[len(tokenList):] {
			OutputRed(t, "not ok ---- Remaining expecting %s '%s'", token.Type.String(), token.Str)
			failure = true
		}
	}
	if failure {
		t.Fatal("See log.")
	}

	return tokens
}

func AssertTokenType(t *testing.T, tokenType ast.TokenType, token *ast.Token) {
	assert.NotNil(t, token)
	if tokenType != token.Type {
		OutputRed(t, "not ok - expecting %s. Got %s '%s'", tokenType.String(), token.Type.String(), token.Str)
	} else {
		OutputGreen(t, "ok - expecting %s. Got %s '%s'", tokenType.String(), token.Type.String(), token.Str)
	}
	assert.Equal(t, tokenType, token.Type)
}

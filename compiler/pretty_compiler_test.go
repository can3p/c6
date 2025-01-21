package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/c9s/c6/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertPrettyCompile(t *testing.T, code string, expected string) {
	var parser = parser.NewParser(nil)
	stmts, err := parser.ParseScss(code)
	require.NoError(t, err)

	var buf bytes.Buffer

	var compiler = NewPrettyCompiler(&buf)

	err = compiler.Compile(parser, stmts)
	require.NoError(t, err)
	assert.Equal(t, expected, strings.TrimSpace(buf.String()))
}

func TestPrettyCompileUniversalSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`* {font-weight: bold; }`,
		`* {
  font-weight: bold;
}`)
}

func TestPrettyCompileClassSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`.foo-bar {font-weight: bold; }`,
		`.foo-bar {
  font-weight: bold;
}`)
}

func TestPrettyCompileIdSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`#myId {font-weight: bold; }`,
		`#myId {
  font-weight: bold;
}`)
}

func TestPrettyCompileAttributeSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`[type=text] {font-weight: bold; }`,
		`[type=text] {
  font-weight: bold;
}`)
}

func TestPrettyCompileAttributeSelectorWithTypeName(t *testing.T) {
	AssertPrettyCompile(t,
		`input[type=text] {font-weight: bold; }`,
		`input[type=text] {
  font-weight: bold;
}`)
}

func TestPrettyCompileSelectorGroup(t *testing.T) {
	AssertPrettyCompile(t,
		`html, span, div {font-weight: bold; }`,
		`html, span, div {
  font-weight: bold;
}`)
}

func TestPrettyCompileCompoundSelector1(t *testing.T) {
	AssertPrettyCompile(t,
		`*.foo.bar {font-weight: bold; }`,
		`*.foo.bar {
  font-weight: bold;
}`)
}

func TestPrettyCompileCompoundSelector2(t *testing.T) {
	AssertPrettyCompile(t,
		`div.foo.bar[href$=pdf] {font-weight: bold; }`,
		`div.foo.bar[href$=pdf] {
  font-weight: bold;
}`)
}

func TestPrettyCompileComplexSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`*.foo.bar > .posts {font-weight: bold; }`,
		`*.foo.bar > .posts {
  font-weight: bold;
}`)
}

func TestPrettyCompileMultipleDeclarations(t *testing.T) {
	AssertPrettyCompile(t,
		`body { font-weight: bold; color: red; }`,
		`body {
  font-weight: bold;
  color: red;
}`)

}

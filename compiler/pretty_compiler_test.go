package compiler

import (
	"bytes"
	"testing"

	"github.com/c9s/c6/parser"
	"github.com/c9s/c6/runtime"
	"github.com/stretchr/testify/assert"
)

func AssertPrettyCompile(t *testing.T, code string, expected string) {
	var context = runtime.NewContext()
	var parser = parser.NewParser(context)
	var stmts = parser.ParseScss(code)
	var buf bytes.Buffer

	var compiler = NewPrettyCompiler(context, &buf)

	err := compiler.Compile(stmts)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}

func TestPrettyCompileUniversalSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`* { }`,
		`* {}`)
}

func TestPrettyCompileClassSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`.foo-bar { }`,
		`.foo-bar {}`)
}

func TestPrettyCompileIdSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`#myId { }`,
		`#myId {}`)
}

func TestPrettyCompileAttributeSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`[type=text] { }`,
		`[type=text] {}`)
}

func TestPrettyCompileAttributeSelectorWithTypeName(t *testing.T) {
	AssertPrettyCompile(t,
		`input[type=text] { }`,
		`input[type=text] {}`)
}

func TestPrettyCompileSelectorGroup(t *testing.T) {
	AssertPrettyCompile(t,
		`html, span, div { }`,
		`html, span, div {}`)
}

func TestPrettyCompileCompoundSelector1(t *testing.T) {
	AssertPrettyCompile(t,
		`*.foo.bar { }`,
		`*.foo.bar {}`)
}

func TestPrettyCompileCompoundSelector2(t *testing.T) {
	AssertPrettyCompile(t,
		`div.foo.bar[href$=pdf] { }`,
		`div.foo.bar[href$=pdf] {}`)
}

func TestPrettyCompileComplexSelector(t *testing.T) {
	AssertPrettyCompile(t,
		`*.foo.bar > .posts { }`,
		`*.foo.bar > .posts {}`)
}

func TestPrettyCompileMultipleDeclarations(t *testing.T) {
	AssertPrettyCompile(t,
		`body { font-weight: bold; color: red; }`,
		`body {
  font-weight: bold;
  color: red;
}`)

}

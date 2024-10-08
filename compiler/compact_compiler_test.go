package compiler

import (
	"testing"

	"github.com/c9s/c6/parser"
	"github.com/c9s/c6/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertCompile(t *testing.T, code string, expected string) {
	var context = runtime.NewContext()
	var parser = parser.NewParser(context)
	stmts, err := parser.ParseScss(code)
	require.NoError(t, err)
	var compiler = NewCompactCompiler(context)
	out, err := compiler.CompileString(stmts)
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestCompilerCompliant(t *testing.T) {
	var context = runtime.NewContext()
	var compiler = NewCompactCompiler(context)
	compiler.EnableCompliant(CSS3Compliant)
	compiler.EnableCompliant(IE7Compliant)
	assert.True(t, compiler.HasCompliant(CSS3Compliant))
	assert.False(t, compiler.HasCompliant(CSS4Compliant))

	assert.True(t, compiler.HasCompliant(IE7Compliant))
	assert.False(t, compiler.HasCompliant(IE8Compliant))
}

func TestCompileUniversalSelector(t *testing.T) {
	AssertCompile(t,
		`* { }`,
		`* {}`)
}

func TestCompileClassSelector(t *testing.T) {
	AssertCompile(t,
		`.foo-bar { }`,
		`.foo-bar {}`)
}

func TestCompileIdSelector(t *testing.T) {
	AssertCompile(t,
		`#myId { }`,
		`#myId {}`)
}

func TestCompileAttributeSelector(t *testing.T) {
	AssertCompile(t,
		`[type=text] { }`,
		`[type=text] {}`)
}

func TestCompileAttributeSelectorWithTypeName(t *testing.T) {
	AssertCompile(t,
		`input[type=text] { }`,
		`input[type=text] {}`)
}

func TestCompileSelectorGroup(t *testing.T) {
	AssertCompile(t,
		`html, span, div { }`,
		`html, span, div {}`)
}

func TestCompileCompoundSelector1(t *testing.T) {
	AssertCompile(t,
		`*.foo.bar { }`,
		`*.foo.bar {}`)
}

func TestCompileCompoundSelector2(t *testing.T) {
	AssertCompile(t,
		`div.foo.bar[href$=pdf] { }`,
		`div.foo.bar[href$=pdf] {}`)
}

func TestCompileComplexSelector(t *testing.T) {
	AssertCompile(t,
		`*.foo.bar > .posts { }`,
		`*.foo.bar > .posts {}`)
}

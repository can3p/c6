package parser

import (
	"os"
	"testing"

	"github.com/c9s/c6/ast"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunParserTest(code string) (*ast.StmtList, error) {
	var p = NewParser(NewContext())
	return p.ParseScss(code)
}

func TestParserGetFileType(t *testing.T) {
	matrix := map[uint]string{
		UnknownFileType: ".css",
		ScssFileType:    ".scss",
		SassFileType:    ".sass",
		EcssFileType:    ".ecss",
	}

	for k, v := range matrix {
		assert.Equal(t, k, getFileTypeByExtension(v))
	}

}

func TestParserParseFile(t *testing.T) {
	testPath := "test/file.scss"
	bs, _ := os.ReadFile(testPath)
	p := NewParser(NewContext())
	_, err := p.ParseFile(os.DirFS("."), testPath)
	if err != nil {
		t.Fatal(err)
	}

	if e := string(bs); e != p.Content {
		t.Fatalf("got: %s wanted: %s", p.Content, e)
	}

	if e := testPath; e != p.File.FileName {
		t.Fatalf("got: %s wanted: %s", p.File.FileName, e)
	}
}

func TestParserEmptyRuleSetWithUniversalSelector(t *testing.T) {
	stmts, err := RunParserTest(`* { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserEmptyRuleSetWithClassSelector(t *testing.T) {
	stmts, err := RunParserTest(`.first-name { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserEmptyRuleSetWithIdSelector(t *testing.T) {
	stmts, err := RunParserTest(`#myId { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserEmptyRuleSetWithTypeSelector(t *testing.T) {
	stmts, err := RunParserTest(`div { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserEmptyRuleSetWithAttributeSelectorAttributeNameOnly(t *testing.T) {
	stmts, err := RunParserTest(`[href] { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserEmptyRuleSetWithAttributeSelectorPrefixMatch(t *testing.T) {
	stmts, err := RunParserTest(`[href^=http] { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserEmptyRuleSetWithAttributeSelectorSuffixMatch(t *testing.T) {
	stmts, err := RunParserTest(`[href$=pdf] { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserEmptyRuleSetWithTypeSelectorGroup(t *testing.T) {
	stmts, err := RunParserTest(`div, span, html { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserEmptyRuleSetWithComplexSelector(t *testing.T) {
	stmts, err := RunParserTest(`div#myId.first-name.last-name, span, html, .first-name, .last-name { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserNestedRuleSetSimple(t *testing.T) {
	stmts, err := RunParserTest(`div, span, html { .foo { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserNestedRuleSetSimple2(t *testing.T) {
	stmts, err := RunParserTest(`div, span, html { .foo { color: red; background: blue; } text-align: text; float: left; }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserNestedRuleWithParentSelector(t *testing.T) {
	stmts, err := RunParserTest(`div, span, html { & { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserPropertyNameBorderWidth(t *testing.T) {
	stmts, err := RunParserTest(`div { border-width: 3px 3px 3px 3px; }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserNestedProperty(t *testing.T) {
	stmts, err := RunParserTest(`div {
		border: {
			width: 3px;
			color: #000;
		}
	}`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserPropertyNameBorderWidthInterpolation(t *testing.T) {
	stmts, err := RunParserTest(`div { border-#{ $width }: 3px 3px 3px 3px; }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserPropertyNameBorderWidthInterpolation2(t *testing.T) {
	stmts, err := RunParserTest(`div { #{ $name }: 3px 3px 3px 3px; }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserPropertyNameBorderWidthInterpolation3(t *testing.T) {
	stmts, err := RunParserTest(`div { #{ $name }-left: 3px; }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserImportRuleWithUnquoteUrl(t *testing.T) {
	stmts, err := RunParserTest(`@import url(../foo.css);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserImportRuleWithUrl(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`@import url("http://foo.com/bar.css");`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

	rule, ok := (stmts.Stmts)[0].(*ast.ImportStmt)
	assert.True(t, ok, "Convert to ImportStmt OK")
	assert.NotNil(t, rule)
}

func TestParserImportRuleWithString(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`@import "foo.css";`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserImportRuleWithMedia(t *testing.T) {
	stmts, err := RunParserTest(`@import url("foo.css") screen;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserImportRuleWithMultipleMediaTypes(t *testing.T) {
	stmts, err := RunParserTest(`@import url("bluish.css") projection, tv;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserImportRuleWithMediaTypeAndColorFeature(t *testing.T) {
	stmts, err := RunParserTest(`@import url(color.css) screen and (color);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserImportRuleWithMediaTypeAndMaxWidthFeature(t *testing.T) {
	stmts, err := RunParserTest(`@import url(color.css) screen and (max-width: 300px);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserImportRuleWithMedia2(t *testing.T) {
	stmts, err := RunParserTest(`@import url("foo.css") screen and (orientation:landscape);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQuerySimple(t *testing.T) {
	stmts, err := RunParserTest(`@media screen { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryNotScreen(t *testing.T) {
	stmts, err := RunParserTest(`@media not screen { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryOnlyScreen(t *testing.T) {
	stmts, err := RunParserTest(`@media only screen { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryAllAndMinWidth(t *testing.T) {
	stmts, err := RunParserTest(`@media all and (min-width:500px) {  .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryMinWidth(t *testing.T) {
	stmts, err := RunParserTest(`@media (min-width:500px) {  .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryOrientationPortrait(t *testing.T) {
	stmts, err := RunParserTest(`@media (orientation: portrait) { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryMultipleWithComma(t *testing.T) {
	stmts, err := RunParserTest(`@media screen and (color), projection and (color) { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryNone(t *testing.T) {
	stmts, err := RunParserTest(`@media { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryNotAndMonoChrome(t *testing.T) {
	stmts, err := RunParserTest(`@media not all and (monochrome) { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryJustAll(t *testing.T) {
	stmts, err := RunParserTest(`@media all { .red { color: red; } }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryWithExpr1(t *testing.T) {
	var code = `
@media #{$media} {
  .sidebar {
    width: 500px;
  }
}
	`
	stmts, err := RunParserTest(code)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryWithExpr2(t *testing.T) {
	var code = `
@media #{$media} and ($feature: $value) {
  .sidebar {
    width: 500px;
  }
}
	`
	stmts, err := RunParserTest(code)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

/*
func TestParserMediaQueryNestedInRuleSet(t *testing.T) {
	var code = `
h6, .h6 {
  margin: 0 0 10px;
  line-height: 20px;
  font-size: 12px; }
  @media screen and (min-width: 960px) {
    h6, .h6 {
      font-size: 13px;
      margin: 0 0 15px;
  }
}
	`
	stmts, err := RunParserTest(code)
	assert.Equal(t, 1, len(stmts.Stmts))
}
*/

func TestParserMediaQueryWithVendorPrefixFeature(t *testing.T) {
	// FIXME: 'min--moz-device-pixel-ratio' will become '-moz-device-pixel-ratio'
	stmts, err := RunParserTest(`@media (-webkit-min-device-pixel-ratio: 2), (min--moz-device-pixel-ratio: 2) {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserMediaQueryNested(t *testing.T) {
	var code = `
@media screen {
  .sidebar {
    @media (orientation: landscape) {
      width: 500px;
    }
  }
}
	`
	stmts, err := RunParserTest(code)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfTrueStmt(t *testing.T) {
	stmts, err := RunParserTest(`@if true {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfFalseElseStmt(t *testing.T) {
	stmts, err := RunParserTest(`@if false {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfFalseOrTrueElseStmt(t *testing.T) {
	stmts, err := RunParserTest(`@if false or true {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfTrueAndTrueOrFalseElseStmt(t *testing.T) {
	stmts, err := RunParserTest(`@if true and true or true {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfTrueAndTrueOrFalseElseStmt2(t *testing.T) {
	stmts, err := RunParserTest(`@if (true and true) or true {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonGreaterThan(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) > 2 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonGreaterEqual(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) >= 2 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonLessThan(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) < 2 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonLessEqual(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) <= 2 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonEqual(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) == 6 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonUnequal(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) != 6 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonUnequalElseIf(t *testing.T) {
	stmts, err := RunParserTest(`@if (3+3) != 6 {  } @else if (3+3) == 6 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfComparisonAndLogicalExpr(t *testing.T) {
	stmts, err := RunParserTest(`@if 3 > 1 and 4 < 10 and 5 > 3 {  } @else if (3+3) == 6 {  } @else {  }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserIfDeclBlock(t *testing.T) {
	_, err := RunParserTest(`
@if $i == 1 {
	color: #111;
} @else if $i == 2 {
	color: #222;
} @else if $i == 3 {
	color: #333;
} @else {
	color: red;
	background: url(../background.png);
}
	`)
	require.NoError(t, err)
}

func TestParserForStmtSimple(t *testing.T) {
	stmts, err := RunParserTest(`@for $var from 1 through 20 { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserForStmtExprReduce(t *testing.T) {
	stmts, err := RunParserTest(`@for $var from 2 * 3 through 20 * 5 + 10 { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserForStmtRangeOperator(t *testing.T) {
	stmts, err := RunParserTest(`@for $var in 1 .. 10 { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserForStmtRangeOperatorWithExpr(t *testing.T) {
	stmts, err := RunParserTest(`@for $var in 2 + 3 .. 10 * 10 { }`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))

}

func TestParserWhileStmt(t *testing.T) {
	code := `
$i: 6;
@while $i > 0 { $i: $i - 2; }
`
	stmts, err := RunParserTest(code)
	require.NoError(t, err)
	assert.Equal(t, 2, len(stmts.Stmts))
}

func TestParserCSS3Gradient(t *testing.T) {
	// some test cases from htmldog
	// @see http://www.htmldog.com/guides/css/advanced/gradients/
	var buffers = []string{
		`div { background: repeating-linear-gradient(white, black 10px, white 20px); }`,
		`div { background: linear-gradient(135deg, hsl(36,100%,50%) 10%, hsl(72,100%,50%) 60%, white 90%); }`,
		`div { background: linear-gradient(black 0, white 100%); }`,
		`div { background: radial-gradient(#06c 0, #fc0 50%, #039 100%); }`,
		`div { background: linear-gradient(red 0%, green 33.3%, blue 66.7%, black 100%); }`,
		`div { background: -webkit-radial-gradient(100px 200px, circle closest-side, black, white); }`,
	}
	for _, buffer := range buffers {
		stmts, err := RunParserTest(buffer)
		require.NoError(t, err)
		assert.Equal(t, 1, len(stmts.Stmts))
	}
}

func TestParserPropertyListExpr(t *testing.T) {
	var buffers []string = []string{
		`div { width: 1px; }`,
		`div { width: 2px 3px; }`,
		`div { width: 4px, 5px, 6px, 7px; }`,
		`div { width: 4px, 5px 6px, 7px; }`,
		`div { width: 10px 3px + 7px 20px; }`,
		// `div { width: 10px, 3px + 7px, 20px; }`,
	}
	for _, buffer := range buffers {
		stmts, err := RunParserTest(buffer)
		require.NoError(t, err)
		assert.Equal(t, 1, len(stmts.Stmts))
	}
}

func TestParserFontCssSlash(t *testing.T) {
	// should be plain CSS, no division
	// TODO: verify this case
	stmts, err := RunParserTest(`.foo { font: 12px/24px; }`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtWithBooleanTrue(t *testing.T) {
	block, err := RunParserTest(`$foo: true;`)
	require.NoError(t, err)
	t.Logf("%+v\n", block)
}

func TestParserAssignStmtWithBooleanFalse(t *testing.T) {
	block, err := RunParserTest(`$foo: false;`)
	require.NoError(t, err)
	t.Logf("%+v\n", block)
}

func TestParserAssignStmtWithNull(t *testing.T) {
	stmts, err := RunParserTest(`$foo: null;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtList(t *testing.T) {
	stmts, err := RunParserTest(`$foo: 1 2 3 4;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtListWithParenthesis(t *testing.T) {
	stmts, err := RunParserTest(`$foo: (1 2 3 4);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtMap(t *testing.T) {
	stmts, err := RunParserTest(`$foo: (bar: 1, foo: 2);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtCommaSepList(t *testing.T) {
	stmts, err := RunParserTest(`$foo: (1,2,3,4);`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtWithMorePlus(t *testing.T) {
	stmts, err := RunParserTest(`$foo: 12px + 20px + 20px;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtWithExprDefaultFlag(t *testing.T) {
	stmts, err := RunParserTest(`$foo: 12px + 20px + 20px !default;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtWithExprOptionalFlag(t *testing.T) {
	block, err := RunParserTest(`$foo: 12px + 20px + 20px !optional;`)
	require.NoError(t, err)
	t.Logf("%+v\n", block)
}

func TestParserAssignStmtWithComplexExpr(t *testing.T) {
	stmts, err := RunParserTest(`$foo: 12px * (20px + 20px) + 4px / 2;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtWithInterpolation(t *testing.T) {
	stmts, err := RunParserTest(`$foo: #{ 10 + 20 }px;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserAssignStmtLengthPlusLength(t *testing.T) {
	stmts, err := RunParserTest(`$foo: 10px + 20px;`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtNumberPlusNumberMulLength(t *testing.T) {
	stmts, err := RunParserTest(`$foo: (10 + 20) * 3px;`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithHexColorAddOperation(t *testing.T) {
	stmts, err := RunParserTest(`$foo: #000 + 10;`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithHexColorMulOperation(t *testing.T) {
	stmts, err := RunParserTest(`$foo: #010101 * 20;`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithHexColorDivOperation(t *testing.T) {
	stmts, err := RunParserTest(`$foo: #121212 / 2;`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithPxValue(t *testing.T) {
	stmts, err := RunParserTest(`$foo: 10px;`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithSolveableVariableRef(t *testing.T) {
	stmts, err := RunParserTest(`
	$a: 10px; 
	$b: 10px;
	$c: $a + $b;
	`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithUnknownVariableRef(t *testing.T) {
	stmts, err := RunParserTest(`
	$a: 10px; 
	$b: 10px;
	$c: 3 * ($a + $b) + $c;
	`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts) // should be 60px + $c
}

func TestParserAssignStmtWithFunctionCall(t *testing.T) {
	stmts, err := RunParserTest(`$foo: go();`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithFunctionCallIntegerArgument(t *testing.T) {
	stmts, err := RunParserTest(`$foo: go(1,2,3);`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithFunctionCallFunctionCallArgument(t *testing.T) {
	stmts, err := RunParserTest(`$foo: go(bar());`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithFunctionCallVariableArgument(t *testing.T) {
	stmts, err := RunParserTest(`$foo: go($a,$b,$c);`)
	require.NoError(t, err)
	t.Logf("%+v\n", stmts)
}

func TestParserMixinSimple(t *testing.T) {
	_, err := RunParserTest(`
@mixin silly-links {
  a {
    color: blue;
    background-color: red;
  }
}
	`)
	require.NoError(t, err)
}

func TestParserMixinArguments(t *testing.T) {
	_, err := RunParserTest(`
@mixin colors($text, $background, $border) {
  color: $text;
  background-color: $background;
  border-color: $border;
}
	`)
	require.NoError(t, err)
}

func TestParserMixinContentDirective(t *testing.T) {
	_, err := RunParserTest(`
@mixin apply-to-ie6-only {
  * html {
    @content;
  }
}
	`)
	require.NoError(t, err)
}

func TestParserExtendClassSelector(t *testing.T) {
	_, err := RunParserTest(`@extend .foo-bar;`)
	require.NoError(t, err)
}

func TestParserExtendIdSelector(t *testing.T) {
	_, err := RunParserTest(`@extend #myId;`)
	require.NoError(t, err)
}

func TestParserExtendComplexSelector(t *testing.T) {
	_, err := RunParserTest(`@extend #myId > .foo-bar;`)
	require.NoError(t, err)
}

func TestParserInclude(t *testing.T) {
	_, err := RunParserTest(`
		@include apply-to-ie6-only;
	`)
	require.NoError(t, err)
}

func TestParserIncludeWithContentBlock(t *testing.T) {
	stmts, err := RunParserTest(`
		@include apply-to-ie6-only {
			color: white;
		};
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionSimple(t *testing.T) {
	stmts, err := RunParserTest(`
@function grid-width($n) {
  @return $n * $grid-width + ($n - 1) * $gutter-width;
}
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionSimple2(t *testing.T) {
	stmts, err := RunParserTest(`
@function exists($name) {
  @return variable-exists($name);
}
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionSimple3(t *testing.T) {
	stmts, err := RunParserTest(`
@function f() { }
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionSimple4(t *testing.T) {
	stmts, err := RunParserTest(`
@function f() {
  $foo: hi;
  @return g();
}
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionSimple5(t *testing.T) {
	stmts, err := RunParserTest(`
@function g() {
  @return variable-exists(foo);
}
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionWithAssignments(t *testing.T) {
	stmts, err := RunParserTest(`
@function g() {
  $a: 2 * 10;
  @return $a * 99;
}
	`)
	require.NoError(t, err)
	assert.Equal(t, 1, len(stmts.Stmts))
}

func TestParserFunctionCallKeywordArguments(t *testing.T) {
	stmts, err := RunParserTest(`
	@function foo($a, $b) {
		@return $a + $b;
	}
	$c: foo($b: 2, $a: 1);
	`)
	require.NoError(t, err)
	assert.Equal(t, 2, len(stmts.Stmts))
}

func TestParserMassiveRules(t *testing.T) {
	var buffers []string = []string{
		`div { width: auto; }`,
		`div { width: 100px }`,
		`div { width: 100pt }`,
		`div { width: 100em }`,
		`div { width: 100rem }`,
		`div { padding: 10px 10px; }`,
		`div { padding: 10px 10px 20px 30px; }`,
		`div { padding: 10px + 10px; }`,
		`div { padding: 10px + 10px * 3; }`,
		`div { color: red; }`,
		`div { color: rgb(255,255,255); }`,
		`div { color: rgba(255,255,255,0); }`,
		`div { background-image: url("../images/foo.png"); }`,
		// `div { color: #ccddee; }`,
	}
	for _, buffer := range buffers {
		t.Logf("%s\n", buffer)
		var p = NewParser(NewContext())
		stmts, err := p.ParseScss(buffer)
		require.NoError(t, err)
		t.Logf("%+v\n", stmts)
	}
}

/*
func TestParserIfStmtTrueCondition(t *testing.T) {
	p := NewParser(NewContext())
	block := p.ParseScss(`
	div {
		@if true {
			color: red;
		}
	}
	`)
	_ = block
}
*/

func TestParserTypeSelector(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`div { width: auto; }`)
	require.NoError(t, err)
	ruleset, ok := (stmts.Stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserClassSelector(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`.foo-bar { width: auto; }`)
	require.NoError(t, err)
	ruleset, ok := (stmts.Stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserDescendantCombinatorSelector(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`
	.foo
	.bar
	.zoo { width: auto; }`)
	require.NoError(t, err)
	ruleset, ok := (stmts.Stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserAdjacentCombinator(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`.foo + .bar { width: auto; }`)
	require.NoError(t, err)
	ruleset, ok := (stmts.Stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserGeneralSiblingCombinator(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`.foo ~ .bar { width: auto; }`)
	require.NoError(t, err)
	ruleset, ok := (stmts.Stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserChildCombinator(t *testing.T) {
	p := NewParser(NewContext())
	stmts, err := p.ParseScss(`.foo > .bar { width: auto; }`)
	require.NoError(t, err)
	ruleset, ok := (stmts.Stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func BenchmarkParserClassSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(NewContext())
		_, _ = p.ParseScss(`.foo-bar {}`)
	}
}

func BenchmarkParserAttributeSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(NewContext())
		_, _ = p.ParseScss(`input[type=text] {  }`)
	}
}

func BenchmarkParserComplexSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(NewContext())
		_, _ = p.ParseScss(`div#myId.first-name.last-name, span, html, .first-name, .last-name { }`)
	}
}

func BenchmarkParserMediaQueryAllAndMinWidth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(NewContext())
		_, _ = p.ParseScss(`@media all and (min-width:500px) {  .red { color: red; } }`)
	}
}

func BenchmarkParserOverAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(NewContext())
		_, _ = p.ParseScss(`div#myId.first-name.last-name {
			.foo-bar {
				color: red;
				background: #fff;
				border-radius: 10px;
			}

			@for $i from 1 through 100 { }
			@if $i == 1 {
				color: #111;
			} @else if $i == 2 {
				color: #222;
			} @else if $i == 3 {
				color: #333;
			} @else {
				color: red;
				background: url(../background.png);
			}

			div { width: auto; }
			div { width: 100px }
			div { width: 100pt }
			div { width: 100em }
			div { width: 100rem }
			div { padding: 10px 10px; }
			div { padding: 10px 10px 20px 30px; }
			div { padding: 10px + 10px; }
			div { padding: 10px + 10px * 3; }
			div { color: red; }
			div { color: rgb(255,255,255); }
			div { color: rgba(255,255,255,0); }
			div { background-image: url("../images/foo.png"); }
		}`)
	}

}

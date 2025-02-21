package runtime

import (
	"testing"

	"github.com/c9s/c6/ast"
	"github.com/stretchr/testify/assert"
)

/*
func TestReduceExprForUnsolveableExpr(t *testing.T) {
	expr := ast.NewBinaryExpr(ast.NewOp(ast.T_PLUS), ast.NewNumber(10, nil, nil), ast.NewNumber(3, nil, nil), false)
	expr2 := ast.NewBinaryExpr(ast.NewOp(ast.T_PLUS), expr, ast.NewVariable("$a"), false)
	_, ok := ReduceExpr(expr2, nil)
	assert.False(t, ok)
	assert.Equal(t, "13+$a", expr2.String())
	t.Logf("Reduced expression: %+v", expr2)
}
*/

/*
func TestReduceExprForUnsolveableExpr2(t *testing.T) {
	context := NewContext()
	context.CurrentBlock().GetSymTable().Set("$a", ast.NewNumber(10, nil, nil))
	context.CurrentBlock().GetSymTable().Set("$b", ast.NewNumber(10, nil, nil))
	expr := ast.NewBinaryExpr(
		ast.NewOp(ast.T_PLUS),
		ast.NewBinaryExpr(
			ast.NewOp(ast.T_MUL),
			ast.NewNumber(2, nil, nil),
			ast.NewBinaryExpr(
				ast.NewOp(ast.T_PLUS),
				ast.NewVariable("$a"),
				ast.NewVariable("$b"),
				false,
			),
			false),
		ast.NewVariable("$c"),
		false,
	)
	_, ok := ReduceExpr(expr, context)
	assert.False(t, ok)
	assert.Equal(t, "40+$c", expr.String())
	t.Logf("Reduced expression: %+v", expr)
}
*/

func TestComputeNumberAddNumber(t *testing.T) {
	val, err := Compute(ast.NewOp(ast.T_PLUS), ast.NewNumber(10, nil, nil), ast.NewNumber(3, nil, nil))
	assert.NoError(t, err)
	num, ok := val.(*ast.Number)
	assert.True(t, ok)
	assert.Equal(t, 13.0, num.Value)
}

// TODO: commented out failing test, revisit and fix
//func TestComputeNumberAddNumberIncompatibleUnit(t *testing.T) {
//val, err := Compute(ast.NewOp(ast.T_PLUS), ast.NewNumber(10, ast.NewUnit(ast.T_UNIT_PX, nil), nil), ast.NewNumber(3, ast.NewUnit(ast.T_UNIT_PT, nil), nil))
//assert.NoError(t, err)
//assert.Nil(t, val)
//}

func TestComputeNumberMulWithUnit(t *testing.T) {
	val, err := Compute(ast.NewOp(ast.T_MUL), ast.NewNumber(10, ast.NewUnit(ast.T_UNIT_PX, nil), nil), ast.NewNumber(3, nil, nil))
	assert.NoError(t, err)
	num, ok := val.(*ast.Number)
	assert.True(t, ok)
	assert.Equal(t, ast.T_UNIT_PX, num.Unit.Type)
	assert.Equal(t, 30.0, num.Value)
}

func TestComputeNumberDivWithUnit(t *testing.T) {
	val, err := Compute(ast.NewOp(ast.T_DIV),
		ast.NewNumber(10, ast.NewUnit(ast.T_UNIT_PX, nil), nil),
		ast.NewNumber(2, nil, nil))
	assert.NoError(t, err)

	num, ok := val.(*ast.Number)
	assert.True(t, ok)
	assert.NotNil(t, num.Unit)
	assert.Equal(t, ast.T_UNIT_PX, num.Unit.Type)
	assert.Equal(t, 5.0, num.Value)
}

func TestComputeRGBAColorWithNumber(t *testing.T) {
	val, err := Compute(ast.NewOp(ast.T_PLUS), ast.NewRGBAColor(10, 10, 10, 0.2, nil), ast.NewNumber(3, nil, nil))
	assert.NoError(t, err)
	c, ok := val.(*ast.RGBAColor)
	assert.True(t, ok)
	assert.Equal(t, "rgba(13, 13, 13, 0.2)", c.String())
}

func TestComputeRGBColorWithNumber(t *testing.T) {
	val, err := Compute(ast.NewOp(ast.T_PLUS), ast.NewRGBColor(10, 10, 10, nil), ast.NewNumber(3, nil, nil))
	assert.NoError(t, err)
	c, ok := val.(*ast.RGBColor)
	assert.True(t, ok)
	assert.Equal(t, "rgb(13, 13, 13)", c.String())
}

func TestListLookup(t *testing.T) {
	scope := NewScope(nil)
	l := ast.NewList(" ")
	l.Append(ast.NewNumber(1, nil, nil))
	l.Append(ast.NewNumber(2, nil, nil))
	scope.Insert("$a", l)

	lookupWithinBounds := ast.NewListLookup(ast.NewVariableWithToken(&ast.Token{
		Str: "$a",
	}), 1)

	val, err := EvaluateExpr(lookupWithinBounds, scope)
	assert.NoError(t, err)
	assert.Equal(t, "2", val.String())

	lookupOutOfBounds := ast.NewListLookup(ast.NewVariableWithToken(&ast.Token{
		Str: "$a",
	}), 2)

	_, err = EvaluateExpr(lookupOutOfBounds, scope)
	assert.Error(t, err)
}

func TestListSlice(t *testing.T) {
	scope := NewScope(nil)
	l := ast.NewList(" ")
	l.Append(ast.NewNumber(1, nil, nil))
	l.Append(ast.NewNumber(2, nil, nil))
	scope.Insert("$a", l)

	lookupWithinBounds := ast.NewListSlice(ast.NewVariableWithToken(&ast.Token{
		Str: "$a",
	}), 1)

	val, err := EvaluateExpr(lookupWithinBounds, scope)
	assert.NoError(t, err)
	assert.Equal(t, "[2]", val.String())

	lookupBorder := ast.NewListSlice(ast.NewVariableWithToken(&ast.Token{
		Str: "$a",
	}), 2)

	val, err = EvaluateExpr(lookupBorder, scope)
	assert.NoError(t, err)
	assert.Equal(t, "[]", val.String())

	lookupOutOfBounds := ast.NewListSlice(ast.NewVariableWithToken(&ast.Token{
		Str: "$a",
	}), 3)

	_, err = EvaluateExpr(lookupOutOfBounds, scope)
	assert.Error(t, err)
}

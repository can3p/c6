package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

/*
Used for Incompatible unit, data type or unsupported operations

TODO: This is not used yet. our compute functions should return error if possible
*/
type ComputeError struct {
	Message string
	Left    ast.Value
	Right   ast.Value
}

func (self ComputeError) Error() string {
	return self.Message
}

/*
Value
*/
type ComputeFunction func(a ast.Value, b ast.Value) ast.Value

const ValueTypeNum = 7

//var computableMatrix [ValueTypeNum][ValueTypeNum]bool = [ValueTypeNum][ValueTypeNum]bool{
//[> NumberValue <]
//{false, false, false, false, false, false, false},

//[> HexColorValue <]
//{false, false, false, false, false, false, false},

//[> RGBAColorValue <]
//{false, false, false, false, false, false, false},

//[> RGBColorValue <]
//{false, false, false, false, false, false, false},
//}

//[>
//*
//Each row: [5]ComputeFunction{ NumberValue, HexColorValue, RGBAColorValue, RGBColorValue }
//*/
//var computeFunctionMatrix [5][5]ComputeFunction = [5][5]ComputeFunction{

//[> NumberValue <]
//{nil, nil, nil, nil, nil},

//[> HexColorValue <]
//{nil, nil, nil, nil, nil},

//[> RGBAColorValue <]
//{nil, nil, nil, nil, nil},

//[> RGBColorValue <]
//{nil, nil, nil, nil, nil},
//}

func Compute(op *ast.Op, a ast.Value, b ast.Value) (ast.Value, error) {
	if op == nil {
		return nil, fmt.Errorf("op can't be nil")
	}
	switch op.Type {

	case ast.T_EQUAL:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {
			case *ast.Boolean:
				return ast.NewBoolean(ta.Value == tb.Value), nil
			}
		case *ast.String:
			switch tb := b.(type) {
			case *ast.String:
				return ast.NewBoolean(ta.Value == tb.Value), nil
			default:
				return nil, fmt.Errorf("Cannot compare string [%s] to non-string [%s]", a, b)
			}
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value == tb.Value), nil
				} else {
					return nil, fmt.Errorf("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_UNEQUAL:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {
			case *ast.Boolean:
				return ast.NewBoolean(ta.Value != tb.Value), nil
			}
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value != tb.Value), nil
				} else {
					return nil, fmt.Errorf("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_GT:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value > tb.Value), nil
				} else {
					return nil, fmt.Errorf("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_GE:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value >= tb.Value), nil
				} else {
					return nil, fmt.Errorf("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_LT:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value < tb.Value), nil
				} else {
					return nil, fmt.Errorf("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_LE:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value <= tb.Value), nil
				} else {
					return nil, fmt.Errorf("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_LOGICAL_AND:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {

			case *ast.Boolean:
				return ast.NewBoolean(ta.Value && tb.Value), nil

			// For other data type, we cast to boolean
			default:
				if bv, ok := b.(ast.BooleanValue); ok {
					return ast.NewBoolean(bv.Boolean()), nil
				}
			}
		}

	case ast.T_LOGICAL_OR:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {

			case *ast.Boolean:
				return ast.NewBoolean(ta.Value || tb.Value), nil

			// For other data type, we cast to boolean
			default:
				if bv, ok := b.(ast.BooleanValue); ok {
					return ast.NewBoolean(bv.Boolean()), nil
				}
			}
		}

	/*
		arith expr
	*/
	case ast.T_PLUS:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberAddNumber(ta, tb)
			case *ast.HexColor:
				return HexColorAddNumber(tb, ta), nil
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorAddNumber(ta, tb), nil
			}
		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorAddNumber(ta, tb), nil
			}
		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorAddNumber(ta, tb), nil
			}
		}
	case ast.T_MINUS:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberSubNumber(ta, tb)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorSubNumber(ta, tb), nil
			}

		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorSubNumber(ta, tb), nil
			}

		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorSubNumber(ta, tb), nil
			}
		}

	case ast.T_DIV:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberDivNumber(ta, tb), nil
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorDivNumber(ta, tb), nil
			}
		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorDivNumber(ta, tb), nil
			}
		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorDivNumber(ta, tb), nil
			}
		}

	case ast.T_MUL:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberMulNumber(ta, tb), nil
			}

		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorMulNumber(ta, tb), nil
			}

		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorMulNumber(ta, tb), nil
			}

		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorMulNumber(ta, tb), nil
			}
		}
	}
	return nil, nil
}

/*
A simple expression means the operands are scalar, and can be evaluated.
*/
func IsSimpleExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if IsValue(e.Left) && IsValue(e.Right) {
			return true
		}
	case *ast.UnaryExpr:
		if IsValue(e.Expr) {
			return true
		}
	}
	return false
}

/*
This function returns true when the val is a scalar value, not an expression.
*/
func IsValue(val ast.Expr) bool {
	switch val.(type) {
	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor, *ast.HSLColor, *ast.HSVColor, *ast.Boolean:
		return true
	}
	return false
}

func EvaluateExprInBooleanContext(anyexpr ast.Expr, scope *Scope) (ast.Value, error) {
	switch expr := anyexpr.(type) {

	case *ast.BinaryExpr:
		return EvaluateBinaryExprInBooleanContext(expr, scope)

	case *ast.UnaryExpr:
		return EvaluateUnaryExprInBooleanContext(expr, scope)

	case *ast.Variable:
		if lval, err := scope.Lookup(expr.NormalizedName()); err != nil {
			return nil, err
		} else {
			return lval, nil
		}

	default:
		if bval, ok := expr.(ast.BooleanValue); ok {
			return ast.NewBoolean(bval.Boolean()), nil
		}
	}
	return nil, nil
}

func EvaluateBinaryExprInBooleanContext(expr *ast.BinaryExpr, scope *Scope) (ast.Value, error) {

	var lval ast.Value
	var rval ast.Value
	var err error

	switch expr := expr.Left.(type) {
	case *ast.UnaryExpr:
		lval, err = EvaluateUnaryExprInBooleanContext(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.BinaryExpr:
		lval, err = EvaluateBinaryExprInBooleanContext(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.Variable:
		if lval, err = scope.Lookup(expr.NormalizedName()); err != nil {
			return nil, err
		}

	default:
		lval = expr
	}

	switch expr := expr.Right.(type) {
	case *ast.UnaryExpr:
		rval, err = EvaluateUnaryExprInBooleanContext(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.BinaryExpr:
		rval, err = EvaluateBinaryExprInBooleanContext(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.Variable:
		if rval, err = scope.Lookup(expr.NormalizedName()); err != nil {
			return nil, err
		}

	default:
		rval = expr
	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil, nil
}

func EvaluateUnaryExprInBooleanContext(expr *ast.UnaryExpr, scope *Scope) (ast.Value, error) {
	var val ast.Value
	var err error

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpr:
		val, err = EvaluateBinaryExpr(t, scope)
		if err != nil {
			return nil, err
		}
	case *ast.UnaryExpr:
		val, err = EvaluateUnaryExpr(t, scope)
		if err != nil {
			return nil, err
		}
	default:
		val = ast.Value(t)
	}

	switch expr.Op.Type {
	case ast.T_LOGICAL_NOT:
		if bval, ok := val.(ast.BooleanValue); ok {
			return ast.NewBoolean(bval.Boolean()), nil
		} else {
			return nil, fmt.Errorf("BooleanValue interface is not support for %v", val)
		}
	}
	return val, nil
}

func EvaluateRGBColor(args *ast.CallArgumentList, scope *Scope) (ast.Value, error) {
	var values []ast.Expr

	if len(args.Args) == 1 {
		v, err := EvaluateExpr(args.Args[0].Value, scope)
		if err != nil {
			return nil, err
		}

		l, ok := v.(*ast.List)
		if !ok {
			return nil, fmt.Errorf("the only argument is not a list")
		}

		if len(l.Exprs) != 3 {
			return nil, fmt.Errorf("rgb color expects 3 arguments but got %s", l.String())
		}

		values = l.Exprs
	} else {
		if len(args.Args) != 3 {
			return nil, fmt.Errorf("rgb color expects 3 arguments but got %s", args.String())
		}

		for _, a := range args.Args {
			v, err := EvaluateExpr(a.Value, scope)
			if err != nil {
				return nil, err
			}

			values = append(values, v)
		}
	}

	ints := [3]uint32{}

	for idx, e := range values {
		num, ok := e.(*ast.Number)
		if !ok {
			return nil, fmt.Errorf("Argument is not a number: %s - %T", e.String(), e)
		}

		rounded := num.Integer()
		if rounded < 0 {
			rounded = 0
		}

		if rounded > 255 {
			rounded = 255
		}

		ints[idx] = uint32(rounded)
	}

	color := ast.NewRGBColor(ints[0], ints[1], ints[2], nil)
	return color, nil
}

func EvaluateHSLColor(args *ast.CallArgumentList, scope *Scope) (ast.Value, error) {
	var values []ast.Expr

	if len(args.Args) == 1 {
		v, err := EvaluateExpr(args.Args[0].Value, scope)
		if err != nil {
			return nil, err
		}

		l, ok := v.(*ast.List)
		if !ok {
			return nil, fmt.Errorf("the only argument is not a list")
		}

		if len(l.Exprs) != 3 {
			return nil, fmt.Errorf("rgb color expects 3 arguments but got %s", l.String())
		}

		values = l.Exprs
	} else {
		if len(args.Args) != 3 {
			return nil, fmt.Errorf("rgb color expects 3 arguments but got %s", args.String())
		}

		for _, a := range args.Args {
			v, err := EvaluateExpr(a.Value, scope)
			if err != nil {
				return nil, err
			}

			values = append(values, v)
		}
	}

	nums := [3]float64{}

	for idx, e := range values {
		num, ok := e.(*ast.Number)
		if !ok {
			return nil, fmt.Errorf("Argument is not a number: %s - %T", e.String(), e)
		}

		fl := num.Double()

		if num.Unit != nil && num.Unit.Type == ast.T_UNIT_PERCENT {
			fl = fl / 100.0
		}

		if fl < 0 {
			fl = 0
		}

		if idx == 0 {
			if fl > 360 {
				fl = 360.0
			}
		} else if fl > 1.0 {
			fl = 1.0
		}

		nums[idx] = fl
	}

	color := ast.NewHSLColor(nums[0], nums[1], nums[2], nil)
	return color, nil
}

func EvaluateFunctionCall(fc *ast.FunctionCall, scope *Scope) (ast.Value, error) {
	// this is lame, we should do better of course
	if fc.Ident.Str == "rgb" {
		return EvaluateRGBColor(fc.Arguments, scope)
	}

	if fc.Ident.Str == "hsl" {
		return EvaluateHSLColor(fc.Arguments, scope)
	}

	// by default we assume that we've encountered a builtin function
	return fc, nil
}

/*
EvaluateExpr calls EvaluateBinaryExpr. except EvaluateExpr
prevents calculate css slash as division.  otherwise it's the same as
EvaluateBinaryExpr.
*/
func EvaluateExpr(expr ast.Expr, scope *Scope) (v ast.Value, err error) {
	//defer func() {
	//fmt.Printf("EvaluateExpr %s %s %T\n", expr, v, expr)
	//}()

	switch t := expr.(type) {

	case *ast.BinaryExpr:
		// For binary expression that is a CSS slash, we evaluate the expression as a literal string (unquoted)
		if t.IsCssSlash() {
			// return string object without quote
			s := ast.NewString(0, t.Left.String()+"/"+t.Right.String(), nil)
			return s, nil
		}
		return EvaluateBinaryExpr(t, scope)

	case *ast.UnaryExpr:
		return EvaluateUnaryExpr(t, scope)

	case *ast.ListLookup:
		if val, err := scope.Lookup(t.Variable.NormalizedName()); err != nil {
			return nil, err
		} else {
			l, ok := val.(*ast.List)

			if !ok {
				return nil, fmt.Errorf("%s is not a list", t.Variable.NormalizedName())
			}

			if t.Idx >= l.Len() {
				return nil, fmt.Errorf("%s lookup is out of bounds, idx = %d, len = %d", t.Variable.NormalizedName(), t.Idx, l.Len())
			}

			return EvaluateExpr(l.Exprs[t.Idx], scope)
		}

	case *ast.ListSlice:
		if val, err := scope.Lookup(t.Variable.NormalizedName()); err != nil {
			return nil, err
		} else {
			l, ok := val.(*ast.List)

			if !ok {
				return nil, fmt.Errorf("%s is not a list", t.Variable.NormalizedName())
			}

			if t.FromIdx > l.Len() {
				return nil, fmt.Errorf("%s lookup is out of bounds, idx = %d, len = %d", t.Variable.NormalizedName(), t.FromIdx, l.Len())
			}

			val := &ast.List{
				Separator: " ",
			}

			if t.FromIdx < l.Len() {
				for _, expr := range l.Exprs[t.FromIdx:] {
					evaluated, err := EvaluateExpr(expr, scope)
					if err != nil {
						return nil, err
					}

					val.Append(evaluated)
				}
			}

			return val, nil
		}

	case *ast.Variable:
		if val, err := scope.Lookup(t.NormalizedName()); err != nil {
			return nil, err
		} else {
			return val, nil
		}

	case *ast.FunctionCall:
		return EvaluateFunctionCall(t, scope)

	case *ast.List:
		val := &ast.List{
			Separator: t.Separator,
		}

		for _, expr := range t.Exprs {
			evaluated, err := EvaluateExpr(expr, scope)
			if err != nil {
				return nil, err
			}

			val.Append(evaluated)
		}

		return val, nil
	default:
		return ast.Value(expr), nil

	}

}

/*
EvaluateBinaryExpr recursively.
*/
func EvaluateBinaryExpr(expr *ast.BinaryExpr, scope *Scope) (ast.Value, error) {
	var lval ast.Value
	var rval ast.Value
	var err error

	switch expr := expr.Left.(type) {

	case *ast.BinaryExpr:
		lval, err = EvaluateBinaryExpr(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.UnaryExpr:
		lval, err = EvaluateUnaryExpr(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.Variable:
		if varVal, err := scope.Lookup(expr.NormalizedName()); err != nil {
			return nil, err
		} else {
			lval = varVal.(ast.Expr)
		}

	default:
		lval = ast.Value(expr)
	}

	switch expr := expr.Right.(type) {

	case *ast.UnaryExpr:
		rval, err = EvaluateUnaryExpr(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.BinaryExpr:
		rval, err = EvaluateBinaryExpr(expr, scope)
		if err != nil {
			return nil, err
		}

	case *ast.Variable:
		if varVal, err := scope.Lookup(expr.NormalizedName()); err != nil {
			return nil, err
		} else {
			rval = varVal.(ast.Expr)
		}

	default:
		rval = ast.Value(expr)
	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil, nil
}

func EvaluateUnaryExpr(expr *ast.UnaryExpr, scope *Scope) (ast.Value, error) {
	var val ast.Value
	var err error

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpr:
		val, err = EvaluateBinaryExpr(t, scope)
		if err != nil {
			return nil, err
		}

	case *ast.UnaryExpr:
		val, err = EvaluateUnaryExpr(t, scope)
		if err != nil {
			return nil, err
		}
	case *ast.Variable:
		if varVal, err := scope.Lookup(t.NormalizedName()); err != nil {
			return nil, err
		} else {
			val = varVal.(ast.Expr)
		}
	default:
		val = ast.Value(t)
	}

	switch expr.Op.Type {
	case ast.T_NOP:
		// do nothing
	case ast.T_LOGICAL_NOT:
		if bVal, ok := val.(ast.BooleanValue); ok {
			val = ast.NewBoolean(bVal.Boolean())
		}
	case ast.T_MINUS:
		switch n := val.(type) {
		case *ast.Number:
			n.Value = -n.Value
		}
	}
	return val, nil
}

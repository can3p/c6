package runtime

import "github.com/c9s/c6/ast"

func CanReduceExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		return !e.IsCssSlash()
	}
	return true
}

/*
Reduce constant expression to constant.

@return (Value, ok)

ok = true means the expression is reduced to simple constant.

The difference between Evaluate*Expr method is:

- `ReduceExpr` returns either value or expression (when there is an unsolved expression)
- `EvaluateBinaryExpr` returns nil if there is an unsolved expression.
*/
func ReduceExpr(expr ast.Expr, context *Context) (ast.Value, bool, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if exprLeft, ok, err := ReduceExpr(e.Left, context); err != nil {
			return nil, false, err
		} else if ok {
			e.Left = exprLeft
		}
		if exprRight, ok, err := ReduceExpr(e.Right, context); err != nil {
			return nil, false, err
		} else if ok {
			e.Right = exprRight
		}

	case *ast.UnaryExpr:

		if retExpr, ok, err := ReduceExpr(e.Expr, context); err != nil {
			return nil, false, err
		} else if ok {
			e.Expr = retExpr
		}

	case *ast.Variable:
		if context == nil {
			return nil, false, nil
		}

		if varVal, ok := context.GetVariable(e.Name); ok {
			return varVal.(ast.Expr), true, nil
		}

	default:
		// it's already an constant value
		return e, true, nil
	}

	if IsSimpleExpr(expr) {
		switch e := expr.(type) {
		case *ast.BinaryExpr:
			rs, err := EvaluateBinaryExpr(e, context)
			if err != nil {
				return nil, false, err
			}
			return rs, true, nil
		case *ast.UnaryExpr:
			rs, err := EvaluateUnaryExpr(e, context)
			if err != nil {
				return nil, false, err
			}
			return rs, true, nil
		}
	}
	// not a constant expression
	return nil, false, nil
}

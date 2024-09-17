package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

const MaxWhileIterations = 10_000

func ExecuteList(scope *Scope, stmts *ast.StmtList) (*ast.StmtList, error) {
	out := &ast.StmtList{}

	for _, stmt := range stmts.Stmts {
		ret, err := ExecuteSingle(scope, stmt)

		if err != nil {
			return nil, err
		}

		out.AppendList(ret)
	}

	return out, nil
}

func ExecuteSingle(scope *Scope, stmt ast.Stmt) (*ast.StmtList, error) {
	switch t := stmt.(type) {
	case *ast.AssignStmt:
		err := executeAssignStmt(scope, t)
		return nil, err
	case *ast.IfStmt:
		return executeIfStmt(scope, t)
	case *ast.ForStmt:
		return executeForStmt(scope, t)
	case *ast.WhileStmt:
		return executeWhileStmt(scope, t)
	case *ast.RuleSet:
		return executeRuleSet(scope, t)
	case *ast.Property:
		return executeProperty(scope, t)
	}

	return nil, fmt.Errorf("Don't know how to execute the statement %v", stmt)
}

func executeWhileStmt(scope *Scope, stmt *ast.WhileStmt) (*ast.StmtList, error) {
	eval := func(e ast.Expr, scope *Scope) (bool, error) {
		v, err := EvaluateExprInBooleanContext(e, scope)

		if err != nil {
			return false, err
		}

		if v == nil {
			return false, nil
		}

		if bval, ok := v.(ast.BooleanValue); ok {
			return bval.Boolean(), nil
		}

		if _, ok := v.(ast.Null); ok {
			return false, nil
		}

		return false, fmt.Errorf("BooleanValue interface is not support for %T", v)
	}

	child := NewScope(scope)
	count := 0
	out := &ast.StmtList{}

	for {
		keepGoing, err := eval(stmt.Condition, child)
		fmt.Println(stmt.Condition, "evals to", keepGoing)

		if err != nil {
			return nil, err
		}

		if !keepGoing {
			break
		}

		l, err := ExecuteList(child, &stmt.Block.Stmts)

		if err != nil {
			return nil, err
		}

		fmt.Println("while body", l)

		out.AppendList(l)
		count++

		if count == MaxWhileIterations {
			return nil, fmt.Errorf("While loop had %d iterations, most probably it's an infinite loop", count)
		}
	}

	return out, nil
}

func executeForStmt(scope *Scope, stmt *ast.ForStmt) (*ast.StmtList, error) {

	eval := func(e ast.Expr) (int, error) {
		v, err := EvaluateExpr(e, scope)

		if err != nil {
			return 0, err
		}

		num, ok := v.(*ast.Number)

		if !ok {
			return 0, fmt.Errorf("Expected to get a number but got %T", v)
		}

		return num.Integer(), nil
	}

	from, err := eval(stmt.From)

	if err != nil {
		return nil, err
	}

	to, err := eval(stmt.To)

	if err != nil {
		return nil, err
	}

	if from == to {
		return nil, nil
	}

	step := 1

	if to < from {
		step = -1
	}

	// to simplify the loop,
	// since we can now treat both loops as exclusive upper limit
	if stmt.Inclusive {
		to = to + step
	}

	out := &ast.StmtList{}

	for from != to {
		child := NewScope(scope)
		child.Insert(stmt.Variable.NormalizedName(), ast.NewNumber(float64(from), nil, nil))

		l, err := ExecuteList(child, &stmt.Block.Stmts)

		if err != nil {
			return nil, err
		}

		out.AppendList(l)
		from = from + step
	}

	return out, nil
}

func executeIfStmt(scope *Scope, stmt *ast.IfStmt) (*ast.StmtList, error) {

	eval := func(e ast.Expr) (bool, error) {
		v, err := EvaluateExprInBooleanContext(e, scope)

		if err != nil {
			return false, err
		}

		if bval, ok := v.(ast.BooleanValue); ok {
			return bval.Boolean(), nil
		}

		return false, fmt.Errorf("BooleanValue interface is not support for %+v", v)
	}

	var out *ast.StmtList

	v, err := eval(stmt.Condition)

	if err != nil {
		return nil, err
	}

	if v {
		out = &stmt.Block.Stmts
	} else {
		for _, elsif := range stmt.ElseIfs {
			v, err := eval(elsif.Condition)

			if err != nil {
				return nil, err
			}

			if v {
				out = &elsif.Block.Stmts
				break
			}
		}

		if out == nil && stmt.ElseBlock != nil {
			out = &stmt.ElseBlock.Stmts
		}
	}

	if out == nil {
		return nil, nil
	}

	child := NewScope(scope)

	return ExecuteList(child, out)
}

func executeAssignStmt(scope *Scope, stmt *ast.AssignStmt) error {
	varName := stmt.Variable.Name
	val, err := EvaluateExpr(stmt.Expr, scope)

	if err != nil {
		return err
	}

	if stmt.Global {
		fmt.Println("Inserting global", varName, val)
		scope.GetGlobal().Insert(varName, val)
	} else {
		fmt.Println("Inserting local", varName, val)
		scope.Insert(varName, val)
	}

	return nil
}

func executeRuleSet(scope *Scope, stmt *ast.RuleSet) (*ast.StmtList, error) {
	child := NewScope(scope)

	res, err := ExecuteList(child, &stmt.Block.Stmts)

	if err != nil {
		return nil, err
	}

	rs := ast.NewRuleSet()
	decl := ast.NewDeclBlock(rs)
	decl.AppendList(res)
	rs.Block = decl
	rs.Selectors = stmt.Selectors

	return &ast.StmtList{
		Stmts: []ast.Stmt{rs},
	}, nil
}

func executeProperty(scope *Scope, stmt *ast.Property) (*ast.StmtList, error) {
	ret := ast.NewProperty(stmt.Name.Token)

	for _, e := range stmt.Values {
		val, err := EvaluateExpr(e, scope)
		if err != nil {
			return nil, err
		}

		ret.Values = append(ret.Values, val)
	}

	return &ast.StmtList{
		Stmts: []ast.Stmt{ret},
	}, nil
}

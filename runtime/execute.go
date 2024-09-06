package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

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
	case *ast.RuleSet:
		return executeRuleSet(scope, t)
	case *ast.Property:
		return executeProperty(scope, t)
	}

	return nil, fmt.Errorf("Don't know how to execute the statement %v", stmt)
}

func executeAssignStmt(scope *Scope, stmt *ast.AssignStmt) error {
	varName := stmt.Variable.Name
	val, err := EvaluateExpr(stmt.Expr, scope)

	if err != nil {
		return err
	}

	scope.Insert(varName, val)

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

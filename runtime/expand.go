package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func ExpandTree(stmts *ast.StmtList) (*ast.StmtList, error) {
	out := &ast.StmtList{}

	for _, stmt := range stmts.Stmts {
		t, ok := stmt.(*ast.RuleSet)

		if !ok {
			return nil, fmt.Errorf("Tree can only contain rule sets, but has variable of type %T", stmt)
		}

		ret, err := expandRuleset(t)

		if err != nil {
			return nil, err
		}

		out.AppendList(ret)
	}

	return out, nil
}

func expandRuleset(rs *ast.RuleSet) (*ast.StmtList, error) {
	// empty ruleset does not emit any output
	if len(rs.Block.Stmts.Stmts) == 0 {
		return nil, nil
	}

	out := &ast.StmtList{}
	collector := []ast.Stmt{}

	for _, stmt := range rs.Block.Stmts.Stmts {
		switch t := stmt.(type) {
		case *ast.Property:
			collector = append(collector, t)
		case *ast.RuleSet:
			if len(collector) > 0 {
				nrs := ast.NewRuleSet()
				nrs.Selectors = rs.Selectors
				bl := ast.NewDeclBlock(nrs)
				bl.AppendList(&ast.StmtList{
					Stmts: collector,
				})
				nrs.Block = bl
				out.Append(nrs)
				collector = []ast.Stmt{}
			}

			expanded, err := expandRuleset(t)

			if err != nil {
				return nil, err
			}

			for _, stmt := range expanded.Stmts {
				t, ok := stmt.(*ast.RuleSet)

				if !ok {
					return nil, fmt.Errorf("Tree can only contain rule sets, but has variable of type %T", stmt)
				}

				childSel := t.Selectors
				parentSel := rs.Selectors
				bl := t.Block
				resultList := &ast.ComplexSelectorList{}

				for _, psel := range *parentSel {
					for _, csel := range *childSel {
						resSel, err := ast.JoinSelectors(psel, csel)
						if err != nil {
							return nil, err
						}

						resultList.Append(resSel)
					}
				}

				nrs := ast.NewRuleSet()
				nrs.Selectors = resultList
				nrs.Block = bl
				out.Append(nrs)
			}
		default:
			return nil, fmt.Errorf("Unexpected node type in the expanded tree: %T", t)
		}
	}

	if len(collector) > 0 {
		nrs := ast.NewRuleSet()
		nrs.Selectors = rs.Selectors
		bl := ast.NewDeclBlock(nrs)
		bl.AppendList(&ast.StmtList{
			Stmts: collector,
		})
		nrs.Block = bl
		out.Append(nrs)
	}

	return out, nil
}

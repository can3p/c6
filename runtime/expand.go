package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func ExpandTree(stmts *ast.StmtList) ([]*ast.StmtList, error) {
	out := []*ast.StmtList{}
	cssImports := &ast.StmtList{}

	for _, stmt := range stmts.Stmts {
		switch t := stmt.(type) {
		case *ast.RuleSet:
			ret, err := expandRuleset(t)

			if err != nil {
				return nil, err
			}

			out = append(out, ret)
		case *ast.CssImportStmt:
			cssImports.Append(t)
		default:
			return nil, fmt.Errorf("Tree can only contain rule sets or css imports, but has variable of type %T", stmt)
		}
	}

	if len(cssImports.Stmts) > 0 {
		if len(out) > 0 {
			// css import should always go in the beginning
			cssImports.AppendList(out[0])
			out[0] = cssImports
		} else {
			// in case there are not rules except imports
			out = append(out, cssImports)
		}
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
				bl := ast.NewDeclBlock()
				bl.AppendList(&ast.StmtList{
					Stmts: collector,
				})
				nrs.Block = bl
				out.Append(nrs)
				collector = []ast.Stmt{}
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

			expanded, err := expandRuleset(nrs)

			if err != nil {
				return nil, err
			}

			out.AppendList(expanded)
		default:
			return nil, fmt.Errorf("Unexpected node type in the expanded tree: %T", t)
		}
	}

	if len(collector) > 0 {
		nrs := ast.NewRuleSet()
		nrs.Selectors = rs.Selectors
		bl := ast.NewDeclBlock()
		bl.AppendList(&ast.StmtList{
			Stmts: collector,
		})
		nrs.Block = bl
		out.Append(nrs)
	}

	return out, nil
}

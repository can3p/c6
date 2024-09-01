package runtime

import "github.com/c9s/c6/ast"

func OptimizeIfStmt(parentBlock *ast.Block, stm *ast.IfStmt) error {

	// TODO: select passed condition and merge block
	// try to simplify the condition without context and symbol table

	//var mergeBlock = false
	//var ignoreBlock = false
	//val, err := EvaluateExprInBooleanContext(stm.Condition, nil)
	//if err != nil {
	//return err
	//}
	//// check if the expression is evaluated
	//if IsValue(val) {
	//if b, ok := val.(ast.BooleanValue); ok {
	//if b.Boolean() {
	//mergeBlock = true
	//} else {
	//ignoreBlock = true
	//}
	//}
	//}

	// TODO: make this code do something
	//if mergeBlock {
	//// TODO: merge with subblock with the current block
	//} else if ignoreBlock {

	//}

	return nil
}

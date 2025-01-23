package ast

type UseStmt struct {
	Stmt
	// Path to the module being used
	Path *String

	// Optional namespace for the module
	Namespace *Ident

	// Optional configuration map
	Config *Map
}

func NewUseStmt() *UseStmt {
	return &UseStmt{}
}

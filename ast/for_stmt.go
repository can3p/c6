package ast

type ForStmt struct {
	Variable  *Variable
	From      Expr
	To        Expr
	Inclusive bool
	Block     *DeclBlock
}

func (stm ForStmt) CanBeStmt() {}

func (stm ForStmt) String() string {
	var inclusive = " inclusive"
	if !stm.Inclusive {
		inclusive = " exclusive"
	}
	return "@for " + stm.Variable.String() + " from " + stm.From.String() + " through " + stm.To.String() + inclusive + " {  }\n"
}

func NewForStmt(variable *Variable) *ForStmt {
	return &ForStmt{
		Variable: variable,
	}
}

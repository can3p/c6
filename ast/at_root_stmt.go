package ast

type AtRootStmt struct {
	Token    *Token
	Block    *DeclBlock
	Selector *ComplexSelector
}

func (stm AtRootStmt) CanBeStmt() {}

func (stm AtRootStmt) String() string {
	return stm.Token.String()
}

func NewAtRootStmtWithToken(tok *Token) *AtRootStmt {
	return &AtRootStmt{
		Token: tok,
	}
}

package ast

import "strings"

type IncludeStmt struct {
	Token        *Token // @include
	MixinIdent   *Token // mixin identitfier
	ArgumentList []Expr
	ContentBlock *DeclBlock // if any
}

// sass spec assumes that $var_name and $var-name mean the same
func (self IncludeStmt) NormalizedName() string {
	return strings.ReplaceAll(self.MixinIdent.Str, "-", "_")
}

func (stm IncludeStmt) CanBeStmt()     {}
func (stm IncludeStmt) String() string { return "IncludeStmt.String()" }

func NewIncludeStmtWithToken(token *Token) *IncludeStmt {
	return &IncludeStmt{
		Token: token,
	}
}

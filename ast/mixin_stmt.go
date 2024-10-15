package ast

import "strings"

// import ""

type MixinStmt struct {
	Token        *Token
	Ident        *Token
	Block        *DeclBlock
	ArgumentList *ArgumentList
}

// sass spec assumes that $var_name and $var-name mean the same
func (self MixinStmt) NormalizedName() string {
	return strings.ReplaceAll(self.Ident.Str, "-", "_")
}

func (stm MixinStmt) CanBeStmt()     {}
func (stm MixinStmt) String() string { return "{mixin}" }

func NewMixinStmtWithToken(tok *Token) *MixinStmt {
	return &MixinStmt{Token: tok}
}

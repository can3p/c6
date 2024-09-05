package ast

import "strings"

type Variable struct {
	Name  string
	Token *Token
}

// sass spec assumes that $var_name and $var-name mean the same
func (self Variable) NormalizedName() string {
	return strings.ReplaceAll(self.Name, "-", "_")
}

func (self Variable) CanBeNode() {}

func (self Variable) String() string {
	return self.Name
}

func NewVariableWithToken(token *Token) *Variable {
	return &Variable{
		Name:  token.Str,
		Token: token,
	}
}

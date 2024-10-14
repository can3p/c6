package ast

type FunctionCall struct {
	Ident     *Token
	Arguments *CallArgumentList
}

func (self FunctionCall) CanBeNode() {}
func (self FunctionCall) String() (out string) {
	return self.Ident.Str + "(" + self.Arguments.String() + ")"
}

func NewFunctionCallWithToken(token *Token) *FunctionCall {
	return &FunctionCall{
		Ident: token,
	}
}

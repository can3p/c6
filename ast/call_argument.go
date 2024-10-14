package ast

type CallArgument struct {
	Value          Expr
	Name           *Token
	VariableLength bool
}

func (arg CallArgument) String() string {
	switch {
	case arg.Name != nil:
		return arg.Name.String() + ": " + arg.Value.String()
	case arg.VariableLength:
		return arg.Value.String() + "..."
	default:
		return arg.Value.String()
	}
}

func NewCallArgumentWithToken(name *Token, v Expr) *CallArgument {
	return &CallArgument{
		Value:          v,
		Name:           name,
		VariableLength: false,
	}
}

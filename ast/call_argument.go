package ast

type CallArgument struct {
	Value          Expr
	Name           *Variable
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

func NewCallArgumentWithToken(name *Variable, v Expr) *CallArgument {
	return &CallArgument{
		Value:          v,
		Name:           name,
		VariableLength: false,
	}
}

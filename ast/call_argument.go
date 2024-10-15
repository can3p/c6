package ast

type CallArgument struct {
	Value          Expr
	Name           *Variable
	VariableLength bool
}

func (arg CallArgument) String() string {
	switch {
	case arg.Name != nil:
		var v = "<no value>"

		if arg.Value != nil {
			v = arg.Value.String()
		}

		return arg.Name.String() + ": " + v
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

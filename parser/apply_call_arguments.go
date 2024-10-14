package parser

import (
	"github.com/c9s/c6/ast"
)

func ApplyCallArguments(protoList *ast.ArgumentList, callList *ast.CallArgumentList) (*ast.CallArgumentList, error) {
	out := &ast.CallArgumentList{}

	for idx, proto := range protoList.Arguments {
		val := callList.Args[idx]
		v := ast.NewVariableWithToken(proto.Name)

		out.Args = append(out.Args, ast.NewCallArgumentWithToken(v, val))
	}

	return out, nil
}

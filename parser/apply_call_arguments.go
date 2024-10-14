package parser

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func ApplyCallArguments(protoList *ast.ArgumentList, callList *ast.CallArgumentList) (*ast.CallArgumentList, error) {
	out := &ast.CallArgumentList{}

	for idx, proto := range protoList.Arguments {
		var val ast.Expr

		if idx < len(callList.Args) {
			val = callList.Args[idx].Value
			fmt.Println("value from args", val)
		} else {
			if proto.DefaultValue != nil {
				val = proto.DefaultValue
				fmt.Println("default value", val)
			} else {
				return nil, fmt.Errorf("argument number %d is not found", idx+1)
			}
		}

		v := ast.NewVariableWithToken(proto.Name)

		fmt.Println(v, val)

		out.Args = append(out.Args, ast.NewCallArgumentWithToken(v, val))
	}

	return out, nil
}

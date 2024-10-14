package parser

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func ApplyCallArguments(protoList *ast.ArgumentList, callList *ast.CallArgumentList) (*ast.CallArgumentList, error) {
	out := &ast.CallArgumentList{}

	kwMap := map[string]*ast.CallArgument{}
	kwArgsIdx := len(callList.Args)

	for idx, arg := range callList.Args {
		if arg.Name != nil {
			if kwArgsIdx == len(callList.Args) {
				kwArgsIdx = idx
			}

			kwMap[arg.Name.NormalizedName()] = arg
			continue
		}

		if kwArgsIdx < len(callList.Args) && arg.Name == nil {
			return nil, fmt.Errorf("Found position argument after a named one")
		}

	}

	for idx, proto := range protoList.Arguments {
		var val ast.Expr

		v := ast.NewVariableWithToken(proto.Name)

		if idx < kwArgsIdx {
			val = callList.Args[idx].Value
		} else if arg, ok := kwMap[v.NormalizedName()]; ok {
			val = arg.Value
		} else {
			if proto.DefaultValue != nil {
				val = proto.DefaultValue
			} else {
				return nil, fmt.Errorf("argument number %d is not found", idx+1)
			}
		}

		out.Args = append(out.Args, ast.NewCallArgumentWithToken(v, val))
	}

	return out, nil
}

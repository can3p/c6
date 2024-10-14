package parser

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func ApplyCallArguments(protoList *ast.ArgumentList, callList *ast.CallArgumentList) (*ast.CallArgumentList, error) {
	out := &ast.CallArgumentList{}

	kwMap := map[string]*ast.CallArgument{}
	kwArgsIdx := len(callList.Args)

	var spreadFound bool

	if l := len(protoList.Arguments); l > 0 && protoList.Arguments[l-1].VariableLength {
		spreadFound = true
	}

	for idx, arg := range callList.Args {
		if arg.Name != nil {
			if spreadFound {
				return nil, fmt.Errorf("Named arguments cannot work with spread operator")
			}

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

		// spread var is always the last one, this is (should be) checked in parser
		if proto.VariableLength {
			list := ast.NewList(" ")
			if idx < len(callList.Args) {
				for _, callArg := range callList.Args[idx:] {
					list.Append(callArg.Value)
				}
			}

			out.Args = append(out.Args, ast.NewCallArgumentWithToken(v, list))
			break
		}

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

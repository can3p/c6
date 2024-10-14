package parser

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

func ApplyCallArguments(protoList *ast.ArgumentList, callList *ast.CallArgumentList) (*ast.CallArgumentList, error) {
	out := &ast.CallArgumentList{}

	if protoList == nil && callList == nil {
		return out, nil
	}

	kwMap := map[string]*ast.CallArgument{}
	kwArgsIdx := len(callList.Args)

	var spreadInProtoFound bool
	var spreadInCallSite *ast.Variable

	if l := len(protoList.Arguments); l > 0 && protoList.Arguments[l-1].VariableLength {
		spreadInProtoFound = true
	}

	if l := len(callList.Args); l > 0 && callList.Args[l-1].VariableLength {
		v, ok := callList.Args[l-1].Value.(*ast.Variable)
		if !ok {
			return nil, fmt.Errorf("only a list could be used with spread operator")
		}

		spreadInCallSite = v
	}

	for idx, arg := range callList.Args {
		if arg.Name != nil {
			if spreadInProtoFound || spreadInCallSite != nil {
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
			if spreadInCallSite != nil && idx >= len(callList.Args)-1 {
				val = ast.NewListSlice(spreadInCallSite, idx-(len(callList.Args)-1))
			} else {

				list := ast.NewList(" ")
				if idx < len(callList.Args) {
					for _, callArg := range callList.Args[idx:] {
						list.Append(callArg.Value)
					}
				}
				val = list
			}

			out.Args = append(out.Args, ast.NewCallArgumentWithToken(v, val))
			break
		}

		if spreadInCallSite != nil && idx >= len(callList.Args)-1 {
			val = ast.NewListLookup(spreadInCallSite, idx-(len(callList.Args)-1))
		} else if idx < kwArgsIdx {
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

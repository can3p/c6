package ast

import "bytes"

type CallArgumentList struct {
	Args []*CallArgument
}

func (arg CallArgumentList) String() string {
	var b bytes.Buffer

	idx := 0
	for _, a := range arg.Args {
		if idx > 0 {
			b.WriteString(", ")
		}
		idx++

		b.WriteString(a.String())
	}

	return b.String()
}

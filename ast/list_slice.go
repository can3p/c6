package ast

import (
	"strconv"
)

type ListSlice struct {
	Variable *Variable
	FromIdx  int
}

func (slice ListSlice) String() string {
	return slice.Variable.String() + "[" + strconv.Itoa(slice.FromIdx) + ":]"
}

func NewListSlice(v *Variable, idx int) *ListSlice {
	return &ListSlice{
		Variable: v,
		FromIdx:  idx,
	}
}

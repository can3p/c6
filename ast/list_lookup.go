package ast

import (
	"strconv"
)

type ListLookup struct {
	Variable *Variable
	Idx      int
}

func (lookup ListLookup) String() string {
	if lookup.Variable == nil {
		return "<no name>[" + strconv.Itoa(lookup.Idx) + "]"
	}

	return lookup.Variable.String() + "[" + strconv.Itoa(lookup.Idx) + "]"
}

func NewListLookup(v *Variable, idx int) *ListLookup {
	return &ListLookup{
		Variable: v,
		Idx:      idx,
	}
}

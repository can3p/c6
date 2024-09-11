package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinSelectors(t *testing.T) {
	var ex = []struct {
		description string
		parent      *ComplexSelector
		child       *ComplexSelector
		result      string
	}{
		{
			description: "simple concatenation",
			parent: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							NewTypeSelector("abc"),
						},
					},
				},
			},
			child: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							NewTypeSelector("def"),
						},
					},
				},
			},
			result: "abc def",
		},
		{
			description: "simple replacement",
			parent: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							NewTypeSelector("abc"),
						},
					},
				},
			},
			child: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							&ParentSelector{},
						},
					},
				},
			},
			result: "abc",
		},
		{
			description: "multiple replacements and complex selector",
			parent: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							NewTypeSelector("abc"),
							NewClassSelector(".pr"),
						},
					},
				},
			},
			child: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							&ParentSelector{},
							NewClassSelector(".cl"),
						},
					},
					{
						Combinator: NewAdjacentCombinatorWithToken(nil),
						CompoundSelector: &CompoundSelector{
							&ParentSelector{},
						},
					},
				},
			},
			result: "abc.pr.cl + abc.pr",
		},
		{
			description: "complex selector with multiple items",
			parent: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							NewTypeSelector("abc"),
							NewClassSelector(".pr"),
						},
					},
					{
						Combinator: NewDescendantCombinator(),
						CompoundSelector: &CompoundSelector{
							NewTypeSelector("def"),
						},
					},
				},
			},
			child: &ComplexSelector{
				ComplexSelectorItems: []*ComplexSelectorItem{
					{
						CompoundSelector: &CompoundSelector{
							&ParentSelector{},
							NewClassSelector(".cl"),
						},
					},
				},
			},
			result: "abc.pr def.cl",
		},
	}

	for _, ex := range ex {
		t.Run(ex.description, func(t *testing.T) {
			joined, err := JoinSelectors(ex.parent, ex.child)

			if assert.NoError(t, err) {
				assert.Equal(t, ex.result, joined.String())
			}
		})
	}
}

package ast

import (
	"fmt"
	"slices"
)

type ComplexSelectorItem struct {
	Combinator       Combinator
	CompoundSelector *CompoundSelector
}

func (item *ComplexSelectorItem) String() (css string) {
	return item.Combinator.String() + item.CompoundSelector.String()
}
func (item *ComplexSelectorItem) CSSString() (css string) { return item.String() }

type ComplexSelector struct {
	CompoundSelector     *CompoundSelector
	ComplexSelectorItems []*ComplexSelectorItem
}

func (self *ComplexSelector) AppendCompoundSelector(comb Combinator, sel *CompoundSelector) {
	self.ComplexSelectorItems = append(self.ComplexSelectorItems, &ComplexSelectorItem{comb, sel})
}

func (self *ComplexSelector) String() (css string) {
	css = self.CompoundSelector.String()
	for _, item := range self.ComplexSelectorItems {
		css += item.Combinator.String()
		css += item.CompoundSelector.String()
	}
	return css
}

func (self *ComplexSelector) Clone() *ComplexSelector {
	compountClone := slices.Clone(*self.CompoundSelector)
	duplicate := &ComplexSelector{
		CompoundSelector: &compountClone,
	}

	for _, item := range self.ComplexSelectorItems {
		compountClone := slices.Clone(*item.CompoundSelector)
		duplicate.ComplexSelectorItems = append(duplicate.ComplexSelectorItems, &ComplexSelectorItem{
			// we're not doing copy of a combinator, since we're not doing anything with it
			Combinator:       item.Combinator,
			CompoundSelector: &compountClone,
		})
	}

	return duplicate
}

func (self *ComplexSelector) CSSString() string { return self.String() }

func NewComplexSelector(sel *CompoundSelector) *ComplexSelector {
	return &ComplexSelector{
		CompoundSelector: sel,
	}
}

// JoinSelectors merges parent and child selector taking into account the parent selector
func JoinSelectors(parent, child *ComplexSelector) (*ComplexSelector, error) {
	// flatten complex selector to process it in the uniform way.
	// we probably need to make it a simple array of ComplexSelectorItem in the
	// first place
	childItems := []*ComplexSelectorItem{
		{
			Combinator:       NewDescendantCombinator(),
			CompoundSelector: child.CompoundSelector,
		},
	}

	childItems = append(childItems, child.ComplexSelectorItems...)

	outItems := []*ComplexSelectorItem{}

	// if parent selector ha snot been found, we should
	// simply prepend parent to child
	parentFound := false

	for _, sel := range childItems {
		// this cannot happen
		if len(*sel.CompoundSelector) == 0 {
			continue
		}

		// parent selector is only allowed at the beginning
		// of a compound selector
		if _, ok := (*sel.CompoundSelector)[0].(*ParentSelector); !ok {
			outItems = append(outItems, sel)
			continue
		}

		parentFound = true

		// if we have found a parent selector we should:
		// - make a copy of parent selector
		parentCopy := parent.Clone()
		// - append remaining compund selector items to the last child of parent
		if len(parentCopy.ComplexSelectorItems) > 0 {

			parentCopy.ComplexSelectorItems[len(parentCopy.ComplexSelectorItems)-1].CompoundSelector.Append((*sel.CompoundSelector)[1:]...)
		} else {
			parentCopy.CompoundSelector.Append((*sel.CompoundSelector)[1:]...)
		}

		// - insert them into out items
		outItems = append(outItems, &ComplexSelectorItem{
			Combinator:       sel.Combinator,
			CompoundSelector: parentCopy.CompoundSelector,
		})
		outItems = append(outItems, parentCopy.ComplexSelectorItems...)
	}

	if !parentFound {
		out := &ComplexSelector{
			CompoundSelector:     parent.CompoundSelector,
			ComplexSelectorItems: []*ComplexSelectorItem{},
		}

		out.ComplexSelectorItems = append(out.ComplexSelectorItems, parent.ComplexSelectorItems...)
		out.ComplexSelectorItems = append(out.ComplexSelectorItems, outItems...)

		return out, nil
	}

	// we need to take the first element from out items and use it as a CompoundSelector field in complex selector
	if len(outItems) == 0 {
		return nil, fmt.Errorf("child selector cannot contain zero items")
	}

	return &ComplexSelector{
		CompoundSelector:     outItems[0].CompoundSelector,
		ComplexSelectorItems: outItems[1:],
	}, nil
}

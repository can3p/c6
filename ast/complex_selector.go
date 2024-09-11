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
	if item.Combinator != nil {
		css += item.Combinator.String()
	}
	if item.CompoundSelector != nil {
		css += item.CompoundSelector.String()
	}

	return css
}
func (item *ComplexSelectorItem) CSSString() (css string) { return item.String() }

type ComplexSelector struct {
	ComplexSelectorItems []*ComplexSelectorItem
}

func (self *ComplexSelector) AppendCompoundSelector(comb Combinator, sel *CompoundSelector) {
	self.ComplexSelectorItems = append(self.ComplexSelectorItems, &ComplexSelectorItem{comb, sel})
}

func (self *ComplexSelector) String() (css string) {
	for _, item := range self.ComplexSelectorItems {
		if item.Combinator != nil {
			css += item.Combinator.String()
		}
		if item.CompoundSelector != nil {
			css += item.CompoundSelector.String()
		}
	}
	return css
}

func (self *ComplexSelector) Clone() *ComplexSelector {
	duplicate := &ComplexSelector{}

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

func NewComplexSelector() *ComplexSelector {
	return &ComplexSelector{
		ComplexSelectorItems: []*ComplexSelectorItem{},
	}
}

// JoinSelectors merges parent and child selector taking into account the parent selector
func JoinSelectors(parent, child *ComplexSelector) (*ComplexSelector, error) {
	// flatten complex selector to process it in the uniform way.
	// we probably need to make it a simple array of ComplexSelectorItem in the
	// first place
	childItems := []*ComplexSelectorItem{}

	childItems = append(childItems, child.ComplexSelectorItems...)

	outItems := []*ComplexSelectorItem{}

	// if parent selector ha snot been found, we should
	// simply prepend parent to child
	parentFound := false

	for _, sel := range childItems {
		if sel.CompoundSelector == nil {
			outItems = append(outItems, sel)
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
		// - append remaining compound selector items to the last child of parent
		parentCopy.ComplexSelectorItems[0].Combinator = sel.Combinator
		parentCopy.ComplexSelectorItems[len(parentCopy.ComplexSelectorItems)-1].CompoundSelector.Append((*sel.CompoundSelector)[1:]...)

		// - insert them into out items
		outItems = append(outItems, parentCopy.ComplexSelectorItems...)
	}

	if !parentFound {
		out := &ComplexSelector{
			ComplexSelectorItems: []*ComplexSelectorItem{},
		}

		var childCombinator Combinator
		l := len(parent.ComplexSelectorItems)
		if parent.ComplexSelectorItems[l-1].CompoundSelector == nil {
			childCombinator = parent.ComplexSelectorItems[l-1].Combinator
			out.ComplexSelectorItems = append(out.ComplexSelectorItems, parent.ComplexSelectorItems[:l-1]...)
		} else {
			out.ComplexSelectorItems = append(out.ComplexSelectorItems, parent.ComplexSelectorItems...)
		}

		// when we're joining with a child, it combinator may not be there (e.g. `div`)
		// or there may be only combinator (e.g. `>`)
		if outItems[0].Combinator == nil {
			firstItem := outItems[0]
			compountClone := slices.Clone(*firstItem.CompoundSelector)

			if childCombinator == nil {
				childCombinator = NewDescendantCombinator()
			}

			out.ComplexSelectorItems = append(out.ComplexSelectorItems, &ComplexSelectorItem{
				Combinator:       childCombinator,
				CompoundSelector: &compountClone,
			})
			out.ComplexSelectorItems = append(out.ComplexSelectorItems, outItems[1:]...)
		} else {
			if childCombinator != nil {
				return nil, fmt.Errorf("parent selector ends with a combinator and child selector starts with one")
			}

			out.ComplexSelectorItems = append(out.ComplexSelectorItems, outItems...)
		}

		return out, nil
	}

	// we need to take the first element from out items and use it as a CompoundSelector field in complex selector
	if len(outItems) == 0 {
		return nil, fmt.Errorf("child selector cannot contain zero items")
	}

	return &ComplexSelector{
		ComplexSelectorItems: outItems,
	}, nil
}

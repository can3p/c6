package ast

/*
This is a SCSS only selector
*/
type ParentSelector struct {
	Token *Token
}

func (self ParentSelector) String() string {
	return "&"
}

func NewParentSelectorWithToken(token *Token) *ParentSelector {
	return &ParentSelector{token}
}

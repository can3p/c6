package ast

type MediaQueryStmt struct {
	MediaQueryList *MediaQueryList
	Block          *DeclBlock
}

type MediaQueryList struct {
	List []*MediaQuery
}

func (list *MediaQueryList) Append(query *MediaQuery) {
	list.List = append(list.List, query)
}

func NewMediaQueryList() *MediaQueryList {
	return &MediaQueryList{}
}

func (stm MediaQueryStmt) CanBeStmt() {}

func NewMediaQueryStmt() *MediaQueryStmt {
	return &MediaQueryStmt{}
}

func (stm MediaQueryStmt) String() (out string) {
	for _, mediaQuery := range stm.MediaQueryList.List {
		out += ", " + mediaQuery.String()
	}
	return out[2:]
}

/*
One MediaQuery may contain media type or media expression.
*/
type MediaQuery struct {
	MediaType *MediaType
	MediaExpr Expr
}

func NewMediaQuery(mediaType *MediaType, expr Expr) *MediaQuery {
	return &MediaQuery{mediaType, expr}
}

func (stm MediaQuery) CSS3String() string {
	return stm.String()
}

func (stm MediaQuery) String() (out string) {
	/*
		{media type} and {media expression}
	*/
	if stm.MediaType != nil {
		out += stm.MediaType.String()
	}
	if stm.MediaExpr != nil {
		if stm.MediaType != nil {
			out += " and "
		}
		out += stm.MediaExpr.String()
	}
	return out
}

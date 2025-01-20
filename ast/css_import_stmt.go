package ast

type Url interface{}

/*
For relative Url:  url(../)
*/
type RelativeUrl string

/*
For url(http:....)
*/
type AbsoluteUrl string

/*
For @import "../string";

Which may present absolute or relative url
*/
type StringUrl string

/*
*
The @import rule syntax is described here:

@see http://www.w3.org/TR/2015/CR-css-cascade-3-20150416/#at-import

hides the style sheet from Netscape 4, IE 3 and 4 (not 4.72)

	@import url(../style.css);

hides the style sheet from Netscape 4, IE 3 and 4 (not 4.72), Konqueror 2, and Amaya 5.1

	@import url("../style.css");

hides the style sheet from Netscape 4, IE 6 and below

	@import url(../style.css) screen;

hides the style sheet from Netscape 4, IE 4 and below, Konqueror 2

	@import "../styles.css";
*/
type CssImportStmt struct {
	Url            Url // if it's wrapped with url(...) or "string"
	MediaQueryList []*MediaQuery
}

func NewCssImportStmt() *CssImportStmt {
	return &CssImportStmt{}
}

func (self CssImportStmt) CanBeStmt() {}

func (self CssImportStmt) String() string { return "CssImportStmt.String()" }

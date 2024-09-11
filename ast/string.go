package ast

import "bytes"

type String struct {
	// Can be `"`, `'` or ``
	Quote byte
	Value string
	Token *Token
}

func (self String) String() string {
	var b bytes.Buffer
	if self.Quote > 0 {
		b.WriteByte(self.Quote)
	}

	b.WriteString(self.Value)

	if self.Quote > 0 {
		b.WriteByte(self.Quote)
	}

	return b.String()
}

/*
Create a string object with quote byte
*/
func NewStringWithQuote(quote byte, token *Token) *String {
	return &String{quote, token.Str, token}
}

func NewStringWithToken(token *Token) *String {
	return &String{0, token.Str, token}
}

func NewString(quote byte, value string, token *Token) *String {
	return &String{quote, value, token}
}

/*
When string length is greater than 0, return true for boolean context
*/
func (str String) Boolean() bool {
	return len(str.Value) > 0
}

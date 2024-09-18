package ast

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelWarn
	LogLevelError
)

type LogStmt struct {
	Directive *Token
	LogLevel  LogLevel
	Expr      Expr
}

func NewLogStmt(e Expr, t *Token, l LogLevel) *LogStmt {
	return &LogStmt{
		Directive: t,
		LogLevel:  l,
		Expr:      e,
	}
}

func (stm LogStmt) CanBeStmt() {}
func (stm LogStmt) String() string {
	return "LogStmt.String()"
}

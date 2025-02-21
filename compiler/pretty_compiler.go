package compiler

import (
	"bytes"
	"fmt"
	"os"

	"github.com/c9s/c6/ast"
	"github.com/c9s/c6/parser"
	"github.com/c9s/c6/runtime"
)

const indentSpace = "  "

var ErrUnknownAstNode = fmt.Errorf("Unknown ast node to compile")

type Option func(p *PrettyCompiler)

func DefaultPrinter(msg any) {
	fmt.Fprintln(os.Stderr, msg)
}

type PrettyCompiler struct {
	Buffer       *bytes.Buffer
	Indent       int
	DebugPrinter runtime.Printer
	WarnPrinter  runtime.Printer
}

func NewPrettyCompiler(buf *bytes.Buffer, o ...Option) *PrettyCompiler {
	c := &PrettyCompiler{
		Buffer:       buf,
		Indent:       0,
		DebugPrinter: DefaultPrinter,
		WarnPrinter:  DefaultPrinter,
	}

	for _, o := range o {
		o(c)
	}

	return c
}

func WithDebug(p runtime.Printer) Option {
	return func(c *PrettyCompiler) {
		c.DebugPrinter = p
	}
}

func WithWarn(p runtime.Printer) Option {
	return func(c *PrettyCompiler) {
		c.WarnPrinter = p
	}
}

func (c *PrettyCompiler) changeIndent(delta int) {
	c.Indent += delta
}

func (c *PrettyCompiler) printLine(l string, printLeadingNewLine bool) {
	if printLeadingNewLine {
		err := c.Buffer.WriteByte('\n')

		if err != nil {
			panic(err)
		}
	}

	i := 0

	for i < c.Indent {
		_, err := c.Buffer.WriteString(indentSpace)

		if err != nil {
			panic(err)
		}
		i++
	}

	c.Buffer.WriteString(l)
}

func (c *PrettyCompiler) printNewline() {
	c.Buffer.WriteByte('\n')
}

func (c *PrettyCompiler) printByte(b byte) {
	c.Buffer.WriteByte(b)
}

func (c *PrettyCompiler) printString(s string) {
	c.Buffer.WriteString(s)
}

func (c *PrettyCompiler) CompileComplexSelectorList(selectorList *ast.ComplexSelectorList) {
	c.printLine(selectorList.String(), false)
}

func (c *PrettyCompiler) CompileValue(v ast.Expr) {
	switch v := v.(type) {
	case *ast.List:
		for idx, expr := range v.Exprs {
			if idx > 0 {
				c.printString(v.Separator)
			}

			c.CompileValue(expr)
		}
	default:
		c.printString(v.String())
	}
}

func (c *PrettyCompiler) CompileDeclBlock(block *ast.DeclBlock) {
	for _, stm := range block.Stmts.Stmts {
		switch stm := stm.(type) {
		case *ast.Property:
			c.printLine(stm.Name.String(), true)
			c.printString(": ")
			for idx, v := range stm.Values {
				if idx > 0 {
					c.printByte(' ')
				}

				c.CompileValue(v)
			}
		default:
			c.printLine(stm.String(), true)
		}
		c.printByte(';')
	}
}

func (c *PrettyCompiler) CompileRuleSet(ruleset *ast.RuleSet) {
	c.CompileComplexSelectorList(ruleset.Selectors)
	c.printLine(" {", false)
	c.changeIndent(1)
	c.CompileDeclBlock(ruleset.Block)
	c.changeIndent(-1)
	c.printLine("}", len(ruleset.Block.Stmts.Stmts) > 0)
}

func (c *PrettyCompiler) CompileExpression(stmt ast.Expr) {
	switch t := stmt.(type) {
	case *ast.BinaryExpr:
		c.CompileExpression(t.Left)
		c.printByte(' ')
		c.printString(t.Op.String())
		c.printByte(' ')
		c.CompileExpression(t.Right)
	case *ast.UnaryExpr:
		c.printString(t.Op.String())
		c.printByte(' ')
		c.CompileExpression(t.Expr)
	case *ast.Token:
		c.printString(t.Str)
	case *ast.MediaFeature:
		c.printByte('(')
		c.CompileExpression(t.Feature)
		if t.Value != nil {
			c.printString(": ")
			c.CompileExpression(t.Value)
		}
		c.printByte(')')
	default:
		fmt.Printf("UKNOWN EXPR TYPE: %T\n", t)
		// TODO: replace String() with the method in this file
		c.printString(t.String())
	}
}

func (c *PrettyCompiler) CompileMediaType(stmt *ast.MediaType) {
	c.CompileExpression(stmt.Expr)
}

func (c *PrettyCompiler) CompileMediaQuery(stmt *ast.MediaQuery) {
	if stmt.MediaType != nil {
		// TODO: replace String() with the method in this file
		c.CompileMediaType(stmt.MediaType)
	}
	if stmt.MediaExpr != nil {
		if stmt.MediaType != nil {
			c.printString(" and ")
		}

		c.CompileExpression(stmt.MediaExpr)
	}
}

func (c *PrettyCompiler) CompileMediaQueryList(stmt *ast.MediaQueryList) {
	if len(stmt.List) == 0 {
		return
	}

	for idx, q := range stmt.List {
		if idx > 0 {
			c.printString(", ")
		}

		c.CompileMediaQuery(q)
	}
}

func (c *PrettyCompiler) CompileCssImport(stmt *ast.CssImportStmt) {
	c.printString(fmt.Sprintf("@import url(%s)", stmt.Url))
	if stmt.MediaQueryList != nil {
		c.printByte(' ')
		c.CompileMediaQueryList(stmt.MediaQueryList)
	}

	c.printByte(';')
}

func (c *PrettyCompiler) CompileStmt(anyStm ast.Stmt) error {
	switch stm := anyStm.(type) {
	case *ast.RuleSet:
		c.CompileRuleSet(stm)
		return nil
	case *ast.CssImportStmt:
		c.CompileCssImport(stm)
		return nil
	case *ast.AssignStmt:
		return nil
	}

	return ErrUnknownAstNode
}

func (c *PrettyCompiler) CompileStmtList(list *ast.StmtList) error {
	for idx, stm := range list.Stmts {
		if idx > 0 {
			c.printNewline()
		}

		if err := c.CompileStmt(stm); err != nil {
			return err
		}
	}

	return nil
}

func (c *PrettyCompiler) CompileRoot(list []*ast.StmtList) error {
	for idx, stm := range list {
		if idx > 0 {
			c.printNewline()
			c.printNewline()
		}

		if err := c.CompileStmtList(stm); err != nil {
			return err
		}
	}

	if c.Buffer.Len() > 0 {
		c.printByte('\n')
	}

	return nil
}

func (c *PrettyCompiler) Compile(gp *parser.GlobalParser, list *ast.StmtList) error {
	scope := runtime.NewScope(nil)
	r := runtime.NewRuntime(gp, c.DebugPrinter, c.WarnPrinter)
	executed, err := r.ExecuteList(scope, list)

	if err != nil {
		return err
	}

	expanded, err := runtime.ExpandTree(executed)

	if err != nil {
		return err
	}

	return c.CompileRoot(expanded)
}

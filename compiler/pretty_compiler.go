package compiler

import (
	"bytes"
	"fmt"
	"os"

	"github.com/c9s/c6/ast"
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

func (c *PrettyCompiler) CompileComplexSelectorList(selectorList *ast.ComplexSelectorList) {
	c.printLine(selectorList.String(), false)
}

func (c *PrettyCompiler) CompileDeclBlock(block *ast.DeclBlock) {
	for _, stm := range block.Stmts.Stmts {
		c.printLine(stm.String(), true)
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

func (c *PrettyCompiler) CompileStmt(anyStm ast.Stmt) error {
	switch stm := anyStm.(type) {
	case *ast.RuleSet:
		c.CompileRuleSet(stm)
		return nil
	case *ast.ImportStmt:
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

func (c *PrettyCompiler) Compile(list *ast.StmtList) error {
	scope := runtime.NewScope(nil)
	r := runtime.NewRuntime(c.DebugPrinter, c.WarnPrinter)
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

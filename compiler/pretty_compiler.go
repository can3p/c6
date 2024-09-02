package compiler

import (
	"bytes"
	"fmt"

	"github.com/c9s/c6/ast"
	"github.com/c9s/c6/runtime"
	"github.com/pkg/errors"
)

const indentSpace = "  "

var ErrUnknownAstNode = fmt.Errorf("Unknown ast node to compile")

type PrettyCompiler struct {
	Context      *runtime.Context
	ContextStack []runtime.Context
	Buffer       *bytes.Buffer
	Indent       int
}

func NewPrettyCompiler(context *runtime.Context, buf *bytes.Buffer) *PrettyCompiler {
	return &PrettyCompiler{
		Context: context,
		Buffer:  buf,
		Indent:  0,
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

		i := 0

		for i < c.Indent {
			_, err := c.Buffer.WriteString(indentSpace)

			if err != nil {
				panic(err)
			}
			i++
		}
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
	if block.Stmts != nil {
		if err := c.Compile(block.Stmts); err != nil {
			panic(err)
		}
	}

	if len(block.SubRuleSets) > 0 {
		for _, r := range block.SubRuleSets {
			c.CompileRuleSet(r)
		}
	}

}

func (c *PrettyCompiler) CompileRuleSet(ruleset *ast.RuleSet) {
	c.Context.PushRuleSet(ruleset)
	defer c.Context.PopRuleSet()

	c.CompileComplexSelectorList(ruleset.Selectors)
	c.printLine(" {", false)
	c.changeIndent(1)
	c.CompileDeclBlock(ruleset.Block)
	c.changeIndent(-1)
	c.printLine("}", len(ruleset.Block.Stmts.Stmts) > 0)
}

func (c *PrettyCompiler) CompileMediaQuery(mq *ast.MediaQuery) {
	c.printLine(mq.String(), false)
}

func (c *PrettyCompiler) CompileMediaQueryStmt(mq *ast.MediaQueryStmt) {
	c.printLine("@media ", true)
	for idx, mediaQuery := range mq.MediaQueryList.List {
		if idx > 0 {
			c.printLine(", ", false)
			c.CompileMediaQuery(mediaQuery)
		}
	}
	c.printLine(mq.String(), false)
	c.printLine(" {", false)
	c.changeIndent(1)
	c.CompileDeclBlock(mq.Block)
	c.changeIndent(-1)
	c.printLine("}", true)
}

func (c *PrettyCompiler) CompileProperty(mq *ast.Property) {
	c.printLine(mq.Name.Name, true)
	c.printLine(": ", false)
	for idx, expr := range mq.Values {
		if idx > 0 {
			c.printByte(' ')
		}

		c.printLine(expr.String(), false)
	}

	c.printByte(';')
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
	case *ast.MediaQueryStmt:
		c.CompileMediaQueryStmt(stm)
		return nil
	case *ast.Property:
		c.CompileProperty(stm)
		return nil
	}

	return errors.Wrapf(ErrUnknownAstNode, "node: %v", anyStm)
}

func (c *PrettyCompiler) Compile(list *ast.StmtList) error {
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

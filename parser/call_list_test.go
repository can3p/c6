package parser

import (
	"fmt"
	"testing"

	"github.com/c9s/c6/ast"
	"github.com/stretchr/testify/require"
)

func TestCallArguments(t *testing.T) {
	var ex = []struct {
		description string
		args        string
		argNum      int
		expected    string
	}{
		{
			description: "no args",
			args:        "",
			argNum:      0,
			expected:    "",
		},
		{
			description: "simple case",
			args:        "1, 2, 3",
			argNum:      3,
			expected:    "1, 2, 3",
		},
		{
			description: "multiple expressions including list",
			args:        "$a, 2 + 3, 4px 5px 6px",
			argNum:      3,
			expected:    "$a, (2+3), [4px 5px 6px]",
		},
		{
			description: "spread operator",
			args:        "$a, 2 + 3, $c...",
			argNum:      3,
			expected:    "$a, (2+3), $c...",
		},
		{
			description: "single named argument",
			args:        "$a, 2 + 3, $c: 2px",
			argNum:      3,
			expected:    "$a, (2+3), $c: 2px",
		},
		{
			description: "multiple named arguments",
			args:        "$a, 2 + 3, $g: 1px 2px 3px, $c: 2px",
			argNum:      4,
			expected:    "$a, (2+3), $g: [1px 2px 3px], $c: 2px",
		},
	}

	for _, ex := range ex {
		t.Run(ex.description, func(t *testing.T) {
			text := fmt.Sprintf(`@include abc(%s);`, ex.args)
			t.Log("test parsing", text)
			stmts, err := RunParserTest(text)
			require.NoError(t, err)

			includeStmt, ok := stmts.Stmts[0].(*ast.IncludeStmt)
			require.True(t, ok)

			require.Equal(t, ex.argNum, len(includeStmt.ArgumentList.Args))
			require.Equal(t, ex.expected, includeStmt.ArgumentList.String())
		})
	}
}

package parser

import (
	"fmt"
	"testing"

	"github.com/c9s/c6/ast"
	"github.com/stretchr/testify/require"
)

func TestApplyCallArguments(t *testing.T) {
	var ex = []struct {
		description string
		args        string
		proto       string
		expected    string
	}{
		{
			description: "simple case",
			args:        "1, 2, 3",
			proto:       "$a, $b, $c",
			expected:    "$a: 1, $b: 2, $c: 3",
		},
	}

	for _, ex := range ex {
		t.Run(ex.description, func(t *testing.T) {
			text := fmt.Sprintf(`@include abc(%s);`, ex.args)
			stmts, err := RunParserTest(text)
			require.NoError(t, err)

			includeStmt, ok := stmts.Stmts[0].(*ast.IncludeStmt)
			require.True(t, ok)

			callArgs := includeStmt.ArgumentList

			text = fmt.Sprintf(`@mixin abc(%s) {};`, ex.proto)
			stmts, err = RunParserTest(text)
			require.NoError(t, err)

			mixinStmt, ok := stmts.Stmts[0].(*ast.MixinStmt)
			require.True(t, ok)

			protoArgs := mixinStmt.ArgumentList

			final, err := ApplyCallArguments(protoArgs, callArgs)
			require.NoError(t, err)

			require.Equal(t, ex.expected, final.String())
		})
	}
}

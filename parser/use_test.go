package parser

import (
	"testing"

	"github.com/c9s/c6/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUseStmt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantPath string
		wantNs   string
		wantErr  bool
	}{
		{
			name:     "basic use",
			input:    `@use "foundation/code";`,
			wantPath: "foundation/code",
			wantNs:   "",
			wantErr:  false,
		},
		{
			name:     "use with namespace",
			input:    `@use "src/corners" as c;`,
			wantPath: "src/corners",
			wantNs:   "c",
			wantErr:  false,
		},
		{
			name:    "use without path",
			input:   `@use;`,
			wantErr: true,
		},
		{
			name:    "use with invalid namespace",
			input:   `@use "src/corners" as;`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(nil)

			stmts, err := parser.ParseScss(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, stmts.Stmts, 1)
			useStmt, ok := stmts.Stmts[0].(*ast.UseStmt)
			assert.True(t, ok)
			assert.Equal(t, tt.wantPath, useStmt.Path.Value)

			if tt.wantNs != "" {
				assert.NotNil(t, useStmt.Namespace)
				assert.Equal(t, tt.wantNs, useStmt.Namespace.Ident)
			} else {
				assert.Nil(t, useStmt.Namespace)
			}
		})
	}
}

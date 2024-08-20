package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/c9s/c6/compiler"
	"github.com/c9s/c6/parser"
	"github.com/c9s/c6/runtime"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "c6",
		Short: "C6 is a very fast SASS compatible compiler",
		Long:  `C6 is a SASS compatible implementation written in Go. But wait! this is not only to implement SASS, but also to improve the language for better consistency, syntax and performance.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fname := args[0]
			var context = runtime.NewContext()
			d := os.DirFS(path.Dir(fname))
			var parser = parser.NewParser(context)
			var stmts, err = parser.ParseFile(d, fname)

			if err != nil {
				return err
			}

			var b bytes.Buffer
			var compiler = compiler.NewPrettyCompiler(context, &b)

			err = compiler.Compile(stmts)

			if err != nil {
				return err
			}

			fmt.Println(b.String())
			return nil
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of C6",
		Long:  `All software has versions. This is C6's`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("C6 SASS Compiler v0.1 -- HEAD")
		},
	}
	rootCmd.AddCommand(versionCmd)

	var compileCmd = &cobra.Command{
		Use:   "compile",
		Short: "Compile some scss from stdin",
		// Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			var context = runtime.NewContext()
			var parser = parser.NewParser(context)
			content, _ := io.ReadAll(os.Stdin)

			var stmts = parser.ParseScss(string(content))
			var b bytes.Buffer
			var compiler = compiler.NewPrettyCompiler(context, &b)

			err := compiler.Compile(stmts)

			if err != nil {
				return err
			}

			fmt.Println(b.String())
			return nil
		},
	}

	compileCmd.Flags().Int("precision", 0, "I don't know the meaning of this flag")
	rootCmd.Flags().Int("precision", 0, "I don't know the meaning of this flag")

	rootCmd.AddCommand(compileCmd)
	rootCmd.Execute()
}

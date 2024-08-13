package main

import (
	"fmt"
	"io"
	"os"

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
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of C6",
		Long:  `All software has versions. This is C6's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("C6 SASS Compiler v0.1 -- HEAD")
		},
	}
	rootCmd.AddCommand(versionCmd)

	var compileCmd = &cobra.Command{
		Use:   "compile",
		Short: "Compile some scss files",
		// Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			var context = runtime.NewContext()
			var parser = parser.NewParser(context)
			content, _ := io.ReadAll(os.Stdin)
			var stmts = parser.ParseScss(string(content))
			var compiler = compiler.NewCompactCompiler(context)
			fmt.Println(compiler.CompileString(stmts))
		},
	}
	rootCmd.AddCommand(compileCmd)
	rootCmd.Execute()
}

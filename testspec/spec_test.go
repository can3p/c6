package testspec

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/c9s/c6/compiler"
	"github.com/c9s/c6/parser"
	"github.com/c9s/c6/runtime"
	"github.com/can3p/go-hrx/hrx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const hrxPath = "../sass-spec/spec"
const failuresList = "failures_list.go"

func writeFailuresListP(names map[string][][]string) {
	var b bytes.Buffer

	b.WriteString("package testspec\n\n")
	b.WriteString("var BlacklistedSpecs = map[string]map[string]string{\n")
	for _, spec := range sortedKeys(names) {
		n := names[spec]
		sort.SliceStable(n, func(i, j int) bool {
			return n[i][0] < n[j][0]
		})

		b.WriteString("\t\"")
		b.WriteString(spec)
		b.WriteString("\": {\n")
		for _, n := range n {
			b.WriteString("\t\t\"")
			b.WriteString(n[0])
			b.WriteString(fmt.Sprintf("\": \"%s\",\n", strings.ReplaceAll(n[1], "\n", " ")))
		}
		b.WriteString("\t},\n")
	}

	b.WriteString("}\n")

	if err := os.WriteFile(failuresList, b.Bytes(), fs.ModePerm); err != nil {
		panic(err)
	}
}

func sortedKeys[A any](m map[string]A) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func getHrxFiles(path string) ([]string, error) {
	names := []string{}

	err := fs.WalkDir(os.DirFS(path), ".", func(p string, d fs.DirEntry, err2 error) error {
		if strings.HasSuffix(d.Name(), ".hrx") {
			names = append(names, p)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return names, nil
}

func getInputFiles(fsys fs.FS) ([]string, error) {
	out := []string{}

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err2 error) error {
		if strings.HasSuffix(p, "input.scss") || strings.HasSuffix(p, "input.sass") {
			out = append(out, p)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func TestSpec(t *testing.T) {
	if os.Getenv("SKIP_SPEC") != "" {
		return
	}

	testFiles, err := getHrxFiles(hrxPath)
	require.NoError(t, err)

	assert.Positive(t, len(testFiles))

	success := 0
	testOnly := os.Getenv("TEST_ONLY")
	generateFailuresList := os.Getenv("GENERATE_FAILURES_LIST") != ""
	ignoreBlacklisted := os.Getenv("IGNORE_BLACKLISTED") != ""
	if testOnly != "" && (generateFailuresList || ignoreBlacklisted) {
		panic("TEST_ONLY and GENERATE_FAILURES_LIST or IGNORE_BLACKLISTED vars can't be set together")
	}
	if ignoreBlacklisted && generateFailuresList {
		panic("TEST_ONLY and GENERATE_FAILURES_LIST or IGNORE_BLACKLISTED vars can't be set together")
	}

	failedSpecs := map[string][][]string{}

	addFailure := func(spec, fname string, reason string) {
		if _, ok := failedSpecs[spec]; !ok {
			failedSpecs[spec] = [][]string{}
		}
		failedSpecs[spec] = append(failedSpecs[spec], []string{fname, reason})
	}

	reportSuccess := func(fname, input string) {
		blacklistedInputs := BlacklistedSpecs[fname]

		if ignoreBlacklisted && blacklistedInputs != nil {
			if _, ok := blacklistedInputs[input]; ok {
				assert.True(t, false, "[example %s] Input: %s - test passes, but is listed in the blacklisted object", fname, input)
				return
			}
		}

		success++
	}

	testedCount := 0
	for _, fname := range testFiles {
		if testOnly != "" && testOnly != fname {
			continue
		}

		blacklistedInputs := BlacklistedSpecs[fname]

		archive, err := hrx.OpenReader(path.Join(hrxPath, fname))
		require.NoErrorf(t, err, "[example %s] open hrx file", fname)

		inputFiles, err := getInputFiles(archive)

		require.NoErrorf(t, err, "[example %s] get input files", fname)

		for _, input := range inputFiles {
			if ignoreBlacklisted && blacklistedInputs != nil {
				if _, ok := blacklistedInputs[input]; ok {
					continue
				}
			}

			t.Logf("Processing Input file: %s - %s", fname, input)

			testedCount++

			if !assert.NotPanicsf(t, func() {
				var context = runtime.NewContext()
				var parser = parser.NewParser(context)

				var stmts, err = parser.ParseFile(archive, input)
				if !assert.NoErrorf(t, err, "[example %s] Input: %s, parse failed", fname, input) {
					addFailure(fname, input, "parse_failure")
					return
				}

				var b bytes.Buffer
				var compiler = compiler.NewPrettyCompiler(context, &b)

				compileErr := compiler.Compile(stmts)

				baseName := path.Dir(input)
				errFname := path.Join(baseName, "error")
				warnFname := path.Join(baseName, "warning")
				outputFname := path.Join(baseName, "output.css")

				if _, err := fs.Stat(archive, warnFname); err == nil {
					if !assert.True(t, false, "[example %s] Input: %s - warning is expected, but we don't handle that yet", fname, input) {
						addFailure(fname, input, "unhandled_warning")
						return
					}
				}

				if _, err := fs.Stat(archive, errFname); err == nil {
					if !assert.Errorf(t, compileErr, "[example %s] Input: %s", fname, input) {
						addFailure(fname, input, "compiler_should_have_errored")
						return
					}

					expected, err := fs.ReadFile(archive, errFname)

					if !assert.NoErrorf(t, err, "[example %s] Input: %s", fname, errFname) ||
						!assert.Equalf(t, string(expected), compileErr.Error(), "[example %s] Input: %s", fname, errFname) {
						addFailure(fname, input, "compiler_error_does_not_match")
						return
					}

					reportSuccess(fname, input)
				} else {
					if !assert.NoErrorf(t, compileErr, "[example %s] Input: %s", fname, input) {
						addFailure(fname, input, fmt.Sprintf("compiler_unexpected_compile_error - %s", compileErr.Error()))
						return
					}

					expected, err := fs.ReadFile(archive, outputFname)
					if !assert.NoErrorf(t, err, "[example %s] Input: %s", fname, input) ||
						!assert.Equalf(t, string(expected), b.String(), "[example %s] Input: %s", fname, input) {
						addFailure(fname, input, "compiler_output_does_not_match")
						return
					}

					reportSuccess(fname, input)
				}
			}, "[example %s] Input: %s", fname, input) {
				addFailure(fname, input, "compiler_panic")
			}
		}
	}

	if generateFailuresList {
		writeFailuresListP(failedSpecs)
	}

	t.Logf("%d/%d spec files were successful", success, testedCount)
}

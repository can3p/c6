package testspec

import (
	"bytes"
	"io/fs"
	"os"
	"path"
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

func writeFailuresListP(names map[string][]string) {
	var b bytes.Buffer

	b.WriteString("package testspec\n\n")
	b.WriteString("var failedSpecs = map[string]map[string]struct{}{\n")
	for spec, n := range names {
		b.WriteString("\t\"")
		b.WriteString(spec)
		b.WriteString("\": {\n")
		for _, n := range n {
			b.WriteString("\t\t\"")
			b.WriteString(n)
			b.WriteString("\": {},\n")
		}
		b.WriteString("\t},\n")
	}

	b.WriteString("}\n")

	if err := os.WriteFile(failuresList, b.Bytes(), fs.ModePerm); err != nil {
		panic(err)
	}
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
	type specStat struct {
		fname    string
		success  bool
		fileStat map[string]bool
	}

	testFiles, err := getHrxFiles(hrxPath)
	require.NoError(t, err)

	assert.Positive(t, len(testFiles))

	totalFiles := len(testFiles)
	success := 0
	stats := []*specStat{}
	testOnly := os.Getenv("TEST_ONLY")
	generateFailuresList := os.Getenv("GENERATE_FAILURES_LIST") != ""
	if testOnly != "" && generateFailuresList {
		panic("TEST_ONLY and GENERATE_FAILURES_LIST vars can't be set together")
	}

	failedSpecs := map[string][]string{}

	addFailure := func(spec, fname string) {
		if _, ok := failedSpecs[spec]; !ok {
			failedSpecs[spec] = []string{}
			failedSpecs[spec] = append(failedSpecs[spec], fname)
		}
	}

	for idx, fname := range testFiles {
		if testOnly != "" && testOnly != fname {
			continue
		}

		st := &specStat{
			fname:    fname,
			success:  true,
			fileStat: map[string]bool{},
		}

		archive, err := hrx.OpenReader(path.Join(hrxPath, fname))
		require.NoErrorf(t, err, "[example %s] open hrx file", fname)

		inputFiles, err := getInputFiles(archive)

		require.NoErrorf(t, err, "[example %s] get input files", fname)

		t.Logf("Processing file %s", fname)
		t.Logf("Input files: %d - %v", idx, inputFiles)

		for _, input := range inputFiles {
			t.Logf("Processing Input file: %s", input)
			assert.NotPanicsf(t, func() {
				var context = runtime.NewContext()
				var parser = parser.NewParser(context)

				var stmts, err = parser.ParseFile(archive, input)
				if !assert.NoErrorf(t, err, "[example %s] Input: %s, parse failed", fname, input) {
					st.fileStat[input] = false
					st.success = false
					stats = append(stats, st)
					addFailure(fname, input)
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
						st.fileStat[input] = false
						st.success = false
						stats = append(stats, st)
						addFailure(fname, input)
						return
					}
				}

				if _, err := fs.Stat(archive, errFname); err == nil {
					if !assert.Errorf(t, compileErr, "[example %s] Input: %s", fname, input) {
						st.fileStat[input] = false
						st.success = false
						stats = append(stats, st)
						return
					}

					expected, err := fs.ReadFile(archive, errFname)

					if !assert.NoErrorf(t, err, "[example %s] Input: %s", fname, errFname) ||
						!assert.Equalf(t, string(expected), compileErr.Error(), "[example %s] Input: %s", fname, errFname) {
						st.fileStat[input] = false
						st.success = false
						stats = append(stats, st)
						addFailure(fname, input)
						return
					}

					st.fileStat[input] = true
					success++
					stats = append(stats, st)

				} else {
					expected, err := fs.ReadFile(archive, outputFname)
					if !assert.NoErrorf(t, err, "[example %s] Input: %s", fname, input) ||
						!assert.Equalf(t, expected, err.Error(), "[example %s] Input: %s", fname, input) {
						st.fileStat[input] = false
						stats = append(stats, st)
						addFailure(fname, input)
						return
					}

					st.fileStat[input] = true
					success++
					stats = append(stats, st)
				}
			}, "[example %s] Input: %s", fname, input)
			break
		}
	}

	if generateFailuresList {
		writeFailuresListP(failedSpecs)
	}

	t.Logf("%d/%d spec files were successful", success, totalFiles)

	require.True(t, false)
}

package testspec

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"slices"
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

		if testOnly == "" && slices.Contains([]int{529, 741, 821, 825, 882, 883, 944, 1265, 1266, 1325, 1339, 1340, 1345, 1370, 1829, 1831, 1832, 1833, 1835, 1837, 1838, 1839, 1853, 1855, 1856, 1857, 1871, 1873, 1874, 1875, 1889, 1891, 1892, 1893, 1907, 1909, 1910, 1911, 1925, 1927, 1928, 1929, 1931, 1935, 1941, 1974, 2118, 2206, 2416}, idx) {
			t.Logf("file %s has an endless loop", fname)
			continue
		}

		t.Logf("Processing file %s", fname)
		t.Logf("Input files: %d - %v", idx, inputFiles)

		for _, input := range inputFiles {
			assert.NotPanicsf(t, func() {
				var context = runtime.NewContext()
				var parser = parser.NewParser(context)

				var stmts, err = parser.ParseFile(archive, input)
				if !assert.NoErrorf(t, err, "[example %s] Input: %s, parse failed", fname, input) {
					st.fileStat[input] = false
					st.success = false
					stats = append(stats, st)
					return
				}

				var b bytes.Buffer
				var compiler = compiler.NewPrettyCompiler(context, &b)

				compileErr := compiler.Compile(stmts)

				baseName := path.Dir(input)
				errFname := path.Join(baseName, "error")
				outputFname := path.Join(baseName, "output.css")

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
						stats = append(stats, st)
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

	t.Logf("%d/%d spec files were successful", success, totalFiles)

	require.True(t, false)
}

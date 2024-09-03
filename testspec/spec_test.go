package testspec

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"testing"

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
			b.WriteString(fmt.Sprintf("\": \"%s\",\n", n[1]))
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

func TestSpec(t *testing.T) {
	if os.Getenv("SKIP_SPEC") != "" {
		return
	}

	testFiles, err := getHrxFiles(hrxPath)
	require.NoError(t, err)

	assert.Positive(t, len(testFiles))

	testOnly := os.Getenv("TEST_ONLY")
	generateFailuresList := os.Getenv("GENERATE_FAILURES_LIST") != ""
	ignoreBlacklisted := os.Getenv("IGNORE_BLACKLISTED") != ""
	if testOnly != "" && (generateFailuresList || ignoreBlacklisted) {
		panic("TEST_ONLY and GENERATE_FAILURES_LIST or IGNORE_BLACKLISTED vars can't be set together")
	}
	if ignoreBlacklisted && generateFailuresList {
		panic("TEST_ONLY and GENERATE_FAILURES_LIST or IGNORE_BLACKLISTED vars can't be set together")
	}

	failedSpecs := RunSpecs(t, testFiles, testOnly, func(fname, input string) bool {
		blacklistedInputs := BlacklistedSpecs[fname]

		if ignoreBlacklisted && blacklistedInputs != nil {
			if _, ok := blacklistedInputs[input]; ok {
				return true
			}
		}

		return false
	})

	if generateFailuresList {
		writeFailuresListP(failedSpecs)
	}
}

package agent

import (
	"strings"
)

func ReviewDiff(diff string) []Finding {
	implFiles := diffedGoFiles(diff)
	if len(implFiles) == 0 {
		return []Finding{passAll()}
	}
	missing := missingTests(implFiles, diff)
	if len(missing) == 0 {
		return []Finding{passAll()}
	}
	findings := make([]Finding, 0, len(missing))
	for _, f := range missing {
		findings = append(findings, Finding{
			Agent:    "testing-quality",
			Verdict:  Block,
			Severity: Critical,
			Finding:  "implementation without test",
			File:     f,
			Fix:      "add " + testFileName(f) + " with a failing test first",
		})
	}
	return findings
}

func diffedGoFiles(diff string) []string {
	var files []string
	for _, line := range strings.Split(diff, "\n") {
		if !strings.HasPrefix(line, "+++ b/") {
			continue
		}
		path := strings.TrimPrefix(line, "+++ b/")
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
	}
	return files
}

func missingTests(implFiles []string, diff string) []string {
	var missing []string
	for _, f := range implFiles {
		if !strings.Contains(diff, "+++ b/"+testFileName(f)) {
			missing = append(missing, f)
		}
	}
	return missing
}

func testFileName(implFile string) string {
	return strings.TrimSuffix(implFile, ".go") + "_test.go"
}

func passAll() Finding {
	return Finding{
		Agent:    "testing-quality",
		Verdict:  Pass,
		Severity: Low,
		Finding:  "all implementation files have corresponding tests",
		File:     "",
		Fix:      "",
	}
}

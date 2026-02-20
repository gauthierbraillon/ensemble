package agent

import (
	"os"
	"strings"
)

func CheckFileWrite(filePath string) Finding {
	if !strings.HasSuffix(filePath, ".go") || strings.HasSuffix(filePath, "_test.go") {
		return Finding{Agent: "testing-quality", Verdict: Pass, Severity: Low, Finding: "not an implementation file"}
	}
	testFile := testFileName(filePath)
	if _, err := os.Stat(testFile); err == nil {
		return Finding{Agent: "testing-quality", Verdict: Pass, Severity: Low, Finding: "test file exists"}
	}
	return Finding{
		Agent:    "testing-quality",
		Verdict:  Block,
		Severity: Critical,
		Finding:  "no test file for " + filePath,
		File:     filePath,
		Fix:      "write " + testFile + " with a failing test first",
	}
}

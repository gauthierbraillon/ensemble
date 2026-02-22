package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gauthierbraillon/ensemble/internal/runner"
)

var exportedSymbol = regexp.MustCompile(`\b(func|type|var|const)\s+[A-Z]\w*`)

func ReviewUX(ctx context.Context, diff string, r runner.Runner) []Finding {
	if !hasExportedAPIChange(diff) {
		return nil
	}
	if r == nil {
		return []Finding{skippedUX("no runner configured")}
	}
	raw, err := r.Run(ctx, uxPrompt(diff))
	if err != nil {
		return []Finding{skippedUX(err.Error())}
	}
	findings, err := parseSWEResponse(raw)
	if err != nil {
		return []Finding{skippedUX(fmt.Sprintf("unparseable response: %s", err))}
	}
	if len(findings) == 0 {
		return []Finding{passedUX()}
	}
	for i := range findings {
		findings[i].Agent = "ux-design"
	}
	return findings
}

func hasExportedAPIChange(diff string) bool {
	inTestFile := false
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+++ b/") {
			inTestFile = strings.HasSuffix(line, "_test.go")
			continue
		}
		if inTestFile {
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
			if exportedSymbol.MatchString(line) {
				return true
			}
		}
	}
	return false
}

func uxPrompt(diff string) string {
	return `You are a UX/API design reviewer. Review the following git diff for exported API surface issues only: naming conventions, Go idiomatic naming, API clarity, consistency with existing patterns.

Do NOT comment on code quality, security, or implementation details — exported API naming and consistency only.

Respond with a JSON array of findings. Each finding must be:
{"agent":"ux-design","verdict":"pass"|"warn"|"block","severity":"low"|"medium"|"high"|"critical","finding":"<one line>","file":"<path:line or empty>","fix":"<one line or empty>"}

If no issues found, respond with exactly: []

Diff:
` + diff
}

func passedUX() Finding {
	return Finding{
		Agent:    "ux-design",
		Verdict:  Pass,
		Severity: Low,
		Finding:  "no API design issues found",
	}
}

func skippedUX(reason string) Finding {
	return Finding{
		Agent:    "ux-design",
		Verdict:  Warn,
		Severity: Low,
		Finding:  "skipped: " + reason,
	}
}

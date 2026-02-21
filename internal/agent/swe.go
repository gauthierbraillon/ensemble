package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gauthierbraillon/ensemble/internal/runner"
)

func ReviewCode(ctx context.Context, diff string, r runner.Runner) []Finding {
	if r == nil {
		return []Finding{skippedSWE("no runner configured")}
	}
	raw, err := r.Run(ctx, swePrompt(diff))
	if err != nil {
		return []Finding{skippedSWE(err.Error())}
	}
	findings, err := parseSWEResponse(raw)
	if err != nil {
		return []Finding{skippedSWE(fmt.Sprintf("unparseable response: %s", err))}
	}
	if len(findings) == 0 {
		return []Finding{passedSWE()}
	}
	for i := range findings {
		findings[i].Agent = "software-engineering"
	}
	return findings
}

func swePrompt(diff string) string {
	return `You are a software engineering reviewer. Review the following git diff for code quality issues only: naming, SOLID principles, duplication, dead code, error handling.

Respond with a JSON array of findings. Each finding must be:
{"agent":"software-engineering","verdict":"pass"|"warn"|"block","severity":"low"|"medium"|"high"|"critical","finding":"<one line>","file":"<path:line or empty>","fix":"<one line or empty>"}

If no issues found, respond with exactly: []

Diff:
` + diff
}

func parseSWEResponse(raw string) ([]Finding, error) {
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start == -1 || end == -1 || end < start {
		return nil, fmt.Errorf("no JSON array found")
	}
	raw = raw[start : end+1]

	var findings []Finding
	if err := json.Unmarshal([]byte(raw), &findings); err != nil {
		return nil, err
	}
	var valid []Finding
	for _, f := range findings {
		if validVerdict(f.Verdict) && validSeverity(f.Severity) {
			valid = append(valid, f)
		}
	}
	return valid, nil
}

func skippedSWE(reason string) Finding {
	return Finding{
		Agent:    "software-engineering",
		Verdict:  Warn,
		Severity: Low,
		Finding:  "skipped: " + reason,
	}
}

func passedSWE() Finding {
	return Finding{
		Agent:    "software-engineering",
		Verdict:  Pass,
		Severity: Low,
		Finding:  "no code quality issues found",
	}
}

func validVerdict(v Verdict) bool {
	return v == Pass || v == Warn || v == Block
}

func validSeverity(s Severity) bool {
	return s == Low || s == Medium || s == High || s == Critical
}

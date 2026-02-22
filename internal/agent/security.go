package agent

import (
	"context"
	"fmt"

	"github.com/gauthierbraillon/ensemble/internal/runner"
)

func ReviewSecurity(ctx context.Context, diff string, r runner.Runner) []Finding {
	if r == nil {
		return []Finding{skippedSecurity("no runner configured")}
	}
	raw, err := r.Run(ctx, securityPrompt(diff))
	if err != nil {
		return []Finding{skippedSecurity(err.Error())}
	}
	findings, err := parseSWEResponse(raw)
	if err != nil {
		return []Finding{skippedSecurity(fmt.Sprintf("unparseable response: %s", err))}
	}
	if len(findings) == 0 {
		return []Finding{passedSecurity()}
	}
	for i := range findings {
		findings[i].Agent = "security"
	}
	return findings
}

func securityPrompt(diff string) string {
	return `You are a security reviewer. Review the following git diff for security issues only: OWASP Top 10, hardcoded secrets, injection patterns (SQL, command, path traversal), insecure deserialization, authentication and authorisation flaws.

Do NOT comment on code quality, naming, or style — security issues only.

Respond with a JSON array of findings. Each finding must be:
{"agent":"security","verdict":"pass"|"warn"|"block","severity":"low"|"medium"|"high"|"critical","finding":"<one line>","file":"<path:line or empty>","fix":"<one line or empty>"}

If no issues found, respond with exactly: []

Diff:
` + diff
}

func skippedSecurity(reason string) Finding {
	return Finding{
		Agent:    "security",
		Verdict:  Warn,
		Severity: Low,
		Finding:  "skipped: " + reason,
	}
}

func passedSecurity() Finding {
	return Finding{
		Agent:    "security",
		Verdict:  Pass,
		Severity: Low,
		Finding:  "no security issues found",
	}
}

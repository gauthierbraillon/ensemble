package agent

type Verdict string
type Severity string

const (
	Pass  Verdict = "pass"
	Warn  Verdict = "warn"
	Block Verdict = "block"

	Low      Severity = "low"
	Medium   Severity = "medium"
	High     Severity = "high"
	Critical Severity = "critical"
)

type Finding struct {
	Agent    string   `json:"agent"`
	Verdict  Verdict  `json:"verdict"`
	Severity Severity `json:"severity"`
	Finding  string   `json:"finding"`
	File     string   `json:"file"`
	Fix      string   `json:"fix"`
}

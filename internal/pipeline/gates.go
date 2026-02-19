// Package pipeline defines the quality gates that every commit must pass.
// Gates run in order. Any failure halts the pipeline (Minimum CD requirement).
package pipeline

// Gate represents a single quality gate in the CD pipeline.
type Gate struct {
	Name        string
	MakeTarget  string
	Description string
}

// Gates returns the ordered list of required quality gates.
// Order is intentional: fast, cheap checks run first to fail early.
func Gates() []Gate {
	return []Gate{
		{"Lint", "lint", "Formatting and style enforcement"},
		{"Type check", "typecheck", "Static type analysis"},
		{"Secret scan", "secrets", "Detect hardcoded secrets"},
		{"SAST", "sast", "Injection pattern detection"},
		{"Build", "build", "Compilation"},
		{"Unit tests", "test", "Unit and sociable unit tests"},
		{"Vuln check", "vulncheck", "Dependency vulnerability scan"},
		{"Contract tests", "test-contracts", "Integration boundary contracts"},
		{"Schema validation", "test-schema", "Migration schema validation"},
	}
}

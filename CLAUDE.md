# k8_ai

AI orchestration CLI tool using multiple Claude Code instances in Kubernetes pods (Minikube).
Enforces good engineering practices by default through a mob-style AI agent team.

## Core Principles

- **TDD first**: no implementation without a failing test. Red → Green → Refactor. No exceptions.
- **ATDD**: acceptance tests are written before code, serve as requirements and living documentation.
- **SOLID, DRY, YAGNI**: enforced at review, not aspirational.
- **Trunk-based development**: no long-lived branches. All work integrates to main daily.
- **CD by default**: every commit on main is deployable. Pipeline is the gate.
- **Lean code and docs**: minimum necessary. No speculative abstractions, no unused variables, no commented-out code.
- **Token discipline**: agents are narrow judges, not broad assistants. Scoped context only.

## Development Workflow (ATDD)

1. Write acceptance test (describes behavior from user/system perspective)
2. Watch it fail
3. Write sociable unit tests for the smallest unit involved
4. Watch them fail
5. Implement minimum code to pass unit tests
6. Verify acceptance test passes
7. Refactor — tests stay green
8. Commit to trunk

Acceptance tests = requirements + documentation. They are never mocked at system boundaries.
Sociable unit tests test through real collaborators (no unnecessary mocking).

## Quality Gates (CI Pipeline)

Every commit to main must pass all gates in order. **No skipping. No --no-verify.**

1. Linting and formatting
2. Static type checking
3. Secret scanning
4. SAST (injection patterns)
5. Compilation / build
6. Unit tests
7. Dependency vulnerability scan
8. Contract tests at every integration boundary
9. Schema migration validation

Pipeline must complete in < 10 minutes (Minimum CD requirement).
If the pipeline is red, all feature work stops until it is green.

## Minimum CD Requirements

- Trunk-based development (main is the only long-lived branch)
- Integrate to trunk minimum daily
- Automated tests run before and after merge
- Main pipeline failure halts all feature work
- Pipeline is the sole deployment method to any environment
- Pipeline determines releasability — its decision is final
- Artifacts are immutable (no post-commit modifications)
- On-demand rollback capability
- Application config deployed alongside artifacts
- Production-like test environments

## Agent Team

7 specialized agents. Each agent is a narrow judge with scoped context.
Agents activate by trigger, not on every change.

| Agent | Trigger | Context scope |
|---|---|---|
| Testing & Quality | Any code change | Changed file diff + existing tests for that module |
| Security | Auth, config, deps, input handling | Changed diff + dependency list + auth-related files |
| Software Engineering | Any code change | Changed diff + immediate imports/callers |
| DevOps / Operations | Dockerfile, CI, infra files | Changed infra files only |
| Continuous Improvement | End of cycle | Session summary (not full transcript) |
| UX / Design | Public API, routes, exported types | API surface only |
| Domain & Product | New story or acceptance criteria | Story text + AC + changed file names |

### Agent output format

Agents return structured findings only. No prose explanation unless severity is high.

```json
{
  "agent": "<role>",
  "verdict": "pass | warn | block",
  "severity": "low | medium | high | critical",
  "finding": "<one line>",
  "file": "<path:line>",
  "fix": "<one line suggestion>"
}
```

`block` verdict halts the cycle. Human (you) is the only one who can override a block.

## Token Optimization Rules

- Agents receive scoped context only (see table above). Never full repo + full history.
- Inter-agent communication is structured JSON, not prose.
- Session state is written to Git, not held in agent context.
- Use Haiku for lint-level checks, Sonnet for code review, Opus for architectural decisions.
- Compress conversation history before passing to any agent. Summary > transcript.
- Maximum context budget per agent: defined per role, enforced by orchestrator.

## What Not To Do

- Do not add features not directly asked for
- Do not mock at system boundaries in acceptance tests
- Do not write implementation before a failing test exists
- Do not skip quality gates for any reason
- Do not create helpers or abstractions for single-use operations
- Do not hold full repo context in any single agent
- Do not use long-lived feature branches
- Do not amend published commits

## Roadmap

### Phase 1 — CD Pipeline (current)
Implement Minimum CD pipeline. This is the foundation everything else runs on.
Definition of done: every commit to main passes all 9 quality gates in < 10 minutes.

### Phase 2 — ATDD Workflow
Implement the acceptance test → unit test → implement cycle as an enforced agent workflow.
Acceptance tests are the source of truth for requirements and documentation.

### Phase 3 — Full Quality Gate Hardening
Strengthen the pipeline with the full gate list defined above.
Each gate has a dedicated agent check mapped to it.

## Tech Stack

- **Language**: Go 1.23+
- **CLI framework**: Cobra + Viper
- **Module**: `github.com/gauthierbraillon/ensemble`
- **Testing**: stdlib `testing` + testify
- **Linting**: golangci-lint v2 (standard defaults + misspell, unconvert, unparam)
- **Formatters**: gofmt + goimports (local prefix enforced)
- **SAST**: gosec
- **Secret scanning**: gitleaks (module: `github.com/zricethezav/gitleaks/v8`)
- **Vuln scan**: govulncheck
- **CI validation**: actionlint
- **Orchestration**: Minikube + Kubernetes (Helm) — Phase 3+
- **Agent base**: Claude Code Task system → Docker Compose → K8s (progressive)
- **CI**: GitHub Actions
- **Quality hooks**: Claude Code Hooks (tdd-guard pattern)
- **Token routing**: Haiku / Sonnet / Opus tiered by task complexity

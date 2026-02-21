# ensemble

CLI orchestrator for AI agent teams. Enforces engineering discipline by default.
You are the conductor. Agents are the ensemble. Quality gates are non-negotiable.

## Core Principles

- **TDD first**: no implementation without a failing test. No exceptions.
- **Lean code**: minimum necessary. No speculative abstractions, no unused variables, no commented-out code, no inline comments. Code must be self-explanatory through naming.
- **Atomic commits**: one commit = one complete RED→GREEN→REFACTOR cycle. Never commit broken state. Never commit partial work.
- **Token discipline**: agents are narrow judges, not broad assistants. Scoped context only. Minimum tokens to make the right call.
- **Trunk-based development**: no long-lived branches. All work integrates to main daily.
- **CD by default**: every commit on main is deployable. Pipeline is the gate.
- **Done means green pipeline**: work is not done until the CD pipeline passes. Local green is necessary but not sufficient.

## Development Workflow (RED → GREEN → REFACTOR → DEPLOY)

ATDD IS TDD — same 4-phase cycle, same test type.

Acceptance tests ARE sociable unit tests. They describe observable behavior through real collaborators.
There is no separate acceptance test layer. Test names express the acceptance criterion.
Contract tests are added only when an external interface boundary exists (HTTP, message queue, etc.).

1. **RED**: write one sociable unit test expressing an acceptance criterion, run it, confirm it fails for the right reason
2. **GREEN**: implement minimum code to pass, run only the failing test
3. **REFACTOR**: clean code and tests, run all tests, stay green
4. **DEPLOY**: atomic commit then run the full pipeline — pipeline is the gate

No mocks of our own code. Mock external systems only (APIs, file system, time).

## Atomic Commits

Each commit must represent one complete, working, tested behaviour change.

- Commit only when all tests are green
- Commit message = the acceptance criterion that was just satisfied
- Use conventional commits: `feat:`, `fix:`, `refactor:`, `test:`, `chore:`
- Never amend published commits
- Never use `--no-verify`

## Token Optimization

Every token sent to an agent has a cost. Waste is a bug.

- Agents receive scoped context only (see Agent Team table). Never full repo + full history.
- Inter-agent communication is structured JSON, not prose.
- Session state is written to Git, not held in agent context.
- Orchestrator passes summaries to agents, never raw transcripts.
- Each agent has a maximum context budget enforced by the orchestrator.
- Diff sent to agents is the minimal relevant chunk, not the full diff.

## Model Tiers

One tier setting applies to all agents. Tier reflects confidence in the process, not task complexity.

| Tier   | Model  | When                                        |
|--------|--------|---------------------------------------------|
| opus   | Opus   | Process still being defined, rules evolving |
| sonnet | Sonnet | Process known, judgment still needed        |
| haiku  | Haiku  | Rules encoded, enforce fast and cheap       |

Bump tier up when something new and unclear emerges. Drop back to haiku when stable.

## Quality Gates (CI Pipeline)

Every commit to main must pass all gates in order. No skipping. No --no-verify.

1. Linting and formatting
2. Static type checking
3. Secret scanning
4. SAST (injection patterns)
5. Compilation / build
6. Unit tests
7. Dependency vulnerability scan
8. Contract tests (only when an external interface boundary exists)
9. Schema migration validation (only when a schema exists)

Pipeline must complete in < 10 minutes (Minimum CD requirement).
If the pipeline is red, all feature work stops until it is green.

## Minimum CD Requirements

Source: minimumcd.org

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

7 specialized agents. Each is a narrow judge with scoped context.
Agents activate by trigger, not on every change.

| Agent                  | Trigger                              | Context scope                                          |
|------------------------|--------------------------------------|--------------------------------------------------------|
| Testing & Quality      | Any code change                      | Changed file diff + existing tests for that module     |
| Security               | Auth, config, deps, input handling   | Changed diff + dependency list + auth-related files    |
| Software Engineering   | Any code change                      | Changed diff + immediate imports/callers               |
| DevOps / Operations    | Dockerfile, CI, infra files          | Changed infra files only                              |
| Continuous Improvement | End of cycle                         | Session summary (not full transcript)                  |
| UX / Design            | Public API, routes, exported types   | API surface only                                       |
| Domain & Product       | New story or acceptance criteria     | Story text + AC + changed file names                   |

### Agent output format

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

`block` verdict halts the cycle. You are the only one who can override a block.

## Architecture

```
you (conductor)
     │
     ▼
orchestrator          ← Claude Code hook fires here automatically
     │
     ├── Testing & Quality   (scoped diff + test files)
     ├── Software Engineering (scoped diff + imports)
     ├── Security            (scoped diff + deps, only if triggered)
     └── ...
     │
     ▼
findings (JSON)
     │
     ├── pass   → silent
     ├── warn   → shown
     └── block  → surfaced to you
```

The CLI feels like Claude Code: a conversation with the group, not commands to individual agents.

## What Not To Do

- Do not add features not directly asked for
- Do not write implementation before a failing test exists
- Do not commit partial or broken work
- Do not skip quality gates for any reason
- Do not create helpers or abstractions for single-use operations
- Do not hold full repo context in any single agent
- Do not pass prose between agents — JSON only
- Do not amend published commits

## Roadmap

### Phase 1 — CD Pipeline ✓
Every commit to main passes all 9 quality gates in < 10 minutes.

### Phase 2 — ATDD Workflow (current)
Orchestrator + hook. You set direction, agents enforce discipline automatically.
Definition of done: Claude Code hook fires ensemble automatically on every file write.

### Phase 3 — Full Agent Team
All 7 agents active. Tier config. Interactive CLI session (chef d'orchestre UX).

## Tech Stack

- **Language**: Go 1.23+
- **CLI framework**: Cobra + Viper
- **Module**: `github.com/gauthierbraillon/ensemble`
- **Testing**: stdlib `testing` + testify
- **Linting**: golangci-lint v2
- **Formatters**: gofmt + goimports
- **SAST**: gosec
- **Secret scanning**: gitleaks
- **Vuln scan**: govulncheck
- **CI**: GitHub Actions
- **Hooks**: Claude Code Hooks (PostToolUse on Write)
- **Model routing**: Haiku / Sonnet / Opus via tier config

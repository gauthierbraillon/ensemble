# ensemble

Claude Code for teams. Same CLI feel, XP/CD-disciplined agent team behind it.

Quality is structural, not advisory. You direct. Agents enforce.

---

## Install

```sh
go install github.com/gauthierbraillon/ensemble@latest
```

`ensemble` lands in `$(go env GOPATH)/bin`. Make sure that is on your `$PATH`.

To install from source:

```sh
git clone https://github.com/gauthierbraillon/ensemble
cd ensemble
go install .
```

## Setup

```sh
ensemble init
```

Creates or updates `.claude/settings.json` in the current directory with the `ensemble` PreToolUse hook. Run once per project. Idempotent — safe to re-run.

## Interactive session

```sh
ensemble
```

Opens a conversation with the agent team. Type naturally. Agents enforce discipline, you decide on blocks.

## Hook (automatic TDD enforcement)

After running `ensemble init`, Claude Code calls `ensemble hook` before every `Write` or `Edit`. Writing `foo.go` without `foo_test.go` on disk hard-blocks Claude Code (exit 2):

```json
{"agent":"testing-quality","verdict":"block","severity":"critical","finding":"no test file for foo.go","file":"foo.go","fix":"write foo_test.go with a failing test first"}
```

## Cycle (post-commit gate)

```sh
git diff HEAD~1 | ensemble cycle
```

Exits 1 if any implementation file in the diff has no corresponding test. Each finding is one JSON line.

## Architecture

```
you
 │
 ▼
orchestrator  ← Claude Code hook fires here automatically
 │
 ├── testing-quality   (TDD enforcement, disk + diff)
 ├── engineering        (SOLID/DRY, Sonnet)
 ├── security           (auth, deps, input, Sonnet)
 └── ...               (4 more agents)
 │
 ▼
pass  → silent
warn  → shown
block → you decide
```

## Model tiers

| Tier   | Model  | When                                 |
|--------|--------|--------------------------------------|
| opus   | Opus   | Process still being defined          |
| sonnet | Sonnet | Process known, judgment still needed |
| haiku  | Haiku  | Rules encoded, enforce fast and cheap|

Set tier via `ENSEMBLE_TIER=sonnet` (default: sonnet).

## Workflow

RED → GREEN → REFACTOR → DEPLOY. No exceptions. See [CLAUDE.md](CLAUDE.md).

## References

- [Minimum CD](https://minimumcd.org)
- [Claude Code Hooks](https://docs.anthropic.com/en/docs/claude-code/hooks)

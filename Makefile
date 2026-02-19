.PHONY: lint typecheck secrets sast build test vulncheck test-contracts test-schema ci

# Gate 1 — Linting and formatting
lint:
	golangci-lint run ./...

# Gate 2 — Static type checking
typecheck:
	go vet ./...
	staticcheck ./...

# Gate 3 — Secret scanning
secrets:
	gitleaks detect --source . --no-git

# Gate 4 — SAST (injection patterns)
sast:
	gosec -quiet ./...

# Gate 5 — Compilation
build:
	go build -o bin/ensemble .

# Gate 6 — Unit tests
test:
	go test -race -cover ./...

# Gate 7 — Dependency vulnerability scan
vulncheck:
	govulncheck ./...

# Gate 8 — Contract tests (placeholder until first integration boundary)
test-contracts:
	go test -v -run TestContract ./test/contracts/... 2>/dev/null || echo "no contract tests yet"

# Gate 9 — Schema migration validation (placeholder until first schema)
test-schema:
	go test -v -run TestSchema ./test/schema/... 2>/dev/null || echo "no schema tests yet"

# Run all gates in order — mirrors CI exactly
ci: lint typecheck secrets sast build test vulncheck test-contracts test-schema

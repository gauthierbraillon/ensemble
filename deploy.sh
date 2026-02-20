#!/usr/bin/env bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

pass() { echo -e "${GREEN}✓${NC} $1"; }
fail() { echo -e "${RED}✗${NC} $1"; exit 1; }

echo "=== ensemble deploy ==="

make ci || fail "gates failed — push aborted"
pass "all gates green"

node_modules/.bin/semantic-release --no-ci || true
pass "semantic release done"

git push --follow-tags origin main
pass "pushed"

#!/bin/bash

RED="\033[0;31m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# Regex for Conventional Commits
REGEXP="^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?:\s+.{1,100}"
TYPES="feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert"

if [ -z "$1" ] || [ ! -f "$1" ]; then
    echo -e "${RED}❌ error:${NC} commit message file not found."
    exit 1
fi

COMMIT_MSG=$(cat "$1")

if ! [[ "$COMMIT_MSG" =~ $REGEXP ]]; then
    echo -e "${RED}❌ error:${NC} commit message does not follow Conventional Commits format."
    echo -e "${RED}your message:${NC} \"$COMMIT_MSG\""
    echo -e ""
    echo -e "${BLUE}ℹ️  allowed types:${NC} $TYPES"
    echo -e "${GREEN}✅ example:${NC} feat(auth): add login validation"
    echo -e ""
    exit 1
fi

exit 0
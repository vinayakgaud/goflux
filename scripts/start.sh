#!/usr/bin/env bash

set -e

BRANCH="$1"

if [ -z "$BRANCH" ]; then
    echo "Usage: task start -- feature/<name>"
    exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
    echo "ERROR: Working tree is not clean."
    echo "Commit or stash your changes before starting new work."
    exit 1
fi

if git show-ref --verify --quiet "refs/heads/$BRANCH"; then
    echo "ERROR: Branch '$BRANCH' already exists."
    exit 1
fi

echo "==> Fetching latest main..."
git fetch origin

echo "==> Switching to main..."
git switch main

echo "==> Updating main..."
git pull --rebase origin main

echo "==> Creating branch: $BRANCH"
git switch -c "$BRANCH"

echo
echo "Ready."
echo "Current branch:"
git branch --show-current
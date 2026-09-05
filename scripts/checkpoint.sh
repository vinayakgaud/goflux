#!/usr/bin/env bash

set -e

echo "==> Formatting..."
task format

echo "==> Running quality checks..."
task check

echo "==> Checking Git diff..."
git diff --check

echo
echo "All checks passed."
echo

read -r -p "Commit message: " commit_message

if [ -z "$commit_message" ]; then
    echo "Commit message cannot be empty."
    exit 1
fi

echo
echo "==> Staging changes..."
git add .

echo "==> Committing..."
git commit -am "$commit_message"

echo
echo "Checkpoint committed successfully."
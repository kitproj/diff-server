package main

import (
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

const commandTimeout = 30 * time.Second

func runCommand(ctx context.Context, dir string, stdout, stderr io.Writer, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	cmdLine := name
	if len(args) > 0 {
		cmdLine += " " + strings.Join(args, " ")
	}
	if name == "bash" && len(args) >= 2 && args[0] == "-c" {
		cmdLine = "bash -c <script>"
	}
	log.Printf("executing: %s (dir=%s)", cmdLine, dir)

	start := time.Now()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	duration := time.Since(start).Round(time.Millisecond)
	if err != nil {
		log.Printf("command finished after %s with error (dir=%s): %v", duration, dir, err)
		return err
	}

	log.Printf("command finished in %s (dir=%s)", duration, dir)
	return nil
}

func gitDiffScript() string {
	return `
REPO_NAME="$1"
DEFAULT_BRANCH=""
if git rev-parse --verify main >/dev/null 2>&1; then
  DEFAULT_BRANCH="main"
elif git rev-parse --verify master >/dev/null 2>&1; then
  DEFAULT_BRANCH="master"
fi

CURRENT_BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null || echo "")

if [ -n "$DEFAULT_BRANCH" ] && [ -n "$CURRENT_BRANCH" ] && [ "$CURRENT_BRANCH" != "$DEFAULT_BRANCH" ]; then
  git diff --src-prefix="a/${REPO_NAME}/" --dst-prefix="b/${REPO_NAME}/" ${DEFAULT_BRANCH}...HEAD
else
  git diff --src-prefix="a/${REPO_NAME}/" --dst-prefix="b/${REPO_NAME}/" HEAD
fi

git ls-files --others --exclude-standard | while IFS= read -r file; do
  if [ -n "$file" ]; then
    git diff --no-index --src-prefix="a/${REPO_NAME}/" --dst-prefix="b/${REPO_NAME}/" /dev/null "$file" 2>/dev/null || true
  fi
done
`
}

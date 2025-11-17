package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed diffs.html
var diffsHTML []byte

func diffsHandler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/x-diff") {
		serveDiffsText(w, r)
	} else {
		serveDiffsHTML(w, r)
	}
}

func serveDiffsHTML(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(diffsHTML); err != nil {
		log.Printf("Failed to write HTML response: %v", err)
	}
}

func findGitRepos(root string) ([]string, error) {
	var repos []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repos = append(repos, repoPath)
			return filepath.SkipDir
		}

		return nil
	})

	return repos, err
}

func serveDiffsText(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	writer := &maxSizeWriter{Writer: w, maxSize: 5 * 1024 * 1024}

	repos, err := findGitRepos(".")
	if err != nil {
		http.Error(w, "Failed to find git repositories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/x-diff; charset=utf-8")

	for _, repoPath := range repos {
		relPath, err := filepath.Rel(".", repoPath)
		if err != nil {
			log.Printf("Failed to get relative path for %s: %v", repoPath, err)
			relPath = repoPath
		}
		if relPath == "." {
			relPath = ""
		}

		repoName := relPath
		if repoName == "" {
			repoName = filepath.Base(repoPath)
		}

		cmd := exec.CommandContext(ctx, "bash", "-c", `
# Determine the default branch (main or master)
DEFAULT_BRANCH=""
if git rev-parse --verify main >/dev/null 2>&1; then
  DEFAULT_BRANCH="main"
elif git rev-parse --verify master >/dev/null 2>&1; then
  DEFAULT_BRANCH="master"
fi

# Get current branch
CURRENT_BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null || echo "")

# Determine what to diff against
if [ -n "$DEFAULT_BRANCH" ] && [ -n "$CURRENT_BRANCH" ] && [ "$CURRENT_BRANCH" != "$DEFAULT_BRANCH" ]; then
  # On a feature branch: show all changes from branch point
  git diff --src-prefix=a/`+repoName+`/ --dst-prefix=b/`+repoName+`/ ${DEFAULT_BRANCH}...HEAD
else
  # On default branch or detached HEAD: show only uncommitted changes
  git diff --src-prefix=a/`+repoName+`/ --dst-prefix=b/`+repoName+`/ HEAD
fi

# Always show untracked files
git ls-files --others --exclude-standard | while IFS= read -r file; do
  if [ -n "$file" ]; then
    git diff --no-index --src-prefix=a/`+repoName+`/ --dst-prefix=b/`+repoName+`/ /dev/null "$file" 2>/dev/null || true
  fi
done
		`)
		cmd.Dir = repoPath
		cmd.Stdout = writer
		cmd.Stderr = writer
		if err := cmd.Run(); err != nil {
			log.Printf("Git command failed for %s: %v", repoPath, err)
		}
	}
}

package main

import (
	"context"
	_ "embed"
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
	_, _ = w.Write(diffsHTML)
}

func serveDiffsText(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	writer := &maxSizeWriter{Writer: w, maxSize: 5 * 1024 * 1024}

	entries, err := os.ReadDir(".")
	if err != nil {
		http.Error(w, "Failed to read workspace directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/x-diff; charset=utf-8")

	for _, entry := range entries {
		repoPath := filepath.Join(".", entry.Name())
		gitDir := filepath.Join(repoPath, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			continue
		}

		repoName := entry.Name()

		cmd := exec.CommandContext(ctx, "bash", "-c", `
git diff --src-prefix=a/`+repoName+`/ --dst-prefix=b/`+repoName+`/ HEAD
git ls-files --others --exclude-standard | while read -r file; do
  git diff --no-index --src-prefix=a/`+repoName+`/ --dst-prefix=b/`+repoName+`/ /dev/null "$file"
done
		`)
		cmd.Dir = repoPath
		cmd.Stdout = writer
		cmd.Stderr = writer
		_ = cmd.Run()
	}
}

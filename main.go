package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed diffs.html
var diffsHTML []byte

var workspaceDir string

type maxSizeWriter struct {
	Writer  io.Writer
	maxSize int
	written int
}

func (w *maxSizeWriter) Write(p []byte) (n int, err error) {
	if w.written+len(p) > w.maxSize {
		remaining := w.maxSize - w.written
		if remaining > 0 {
			n, err = w.Writer.Write(p[:remaining])
			w.written += n
		}
		return n, io.ErrShortWrite
	}
	n, err = w.Writer.Write(p)
	w.written += n
	return n, err
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.StringVar(&workspaceDir, "C", "", "Directory to scan for git repositories")
	flag.Parse()

	if workspaceDir == "" {
		var err error
		workspaceDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get current directory: %v\n", err)
			os.Exit(1)
		}
	}

	http.HandleFunc("/", diffsHandler)

	fmt.Printf("Starting server on http://localhost:%s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}

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

	entries, err := os.ReadDir(workspaceDir)
	if err != nil {
		http.Error(w, "Failed to read workspace directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/x-diff; charset=utf-8")

	for _, entry := range entries {
		repoPath := filepath.Join(workspaceDir, entry.Name())
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

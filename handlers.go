package main

import (
	"context"
	"encoding/json"
	_ "embed"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:embed diffs.html
var diffsHTML []byte

func diffsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/repos":
		serveRepos(w, r)
	case "/diff":
		serveRepoDiff(w, r)
	case "/":
		if strings.Contains(r.Header.Get("Accept"), "text/x-diff") {
			serveDiffsText(w, r)
		} else {
			serveDiffsHTML(w, r)
		}
	default:
		http.NotFound(w, r)
	}
}

func serveDiffsHTML(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(diffsHTML); err != nil {
		log.Printf("Failed to write HTML response: %v", err)
	}
}

func serveRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := findGitRepos(".")
	if err != nil {
		http.Error(w, "Failed to find git repositories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	names := make([]string, 0, len(repos))
	for _, repoPath := range repos {
		names = append(names, repoDisplayName(repoPath))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(names); err != nil {
		log.Printf("Failed to encode repos response: %v", err)
	}
}

func serveRepoDiff(w http.ResponseWriter, r *http.Request) {
	repoName := r.URL.Query().Get("repo")
	if repoName == "" {
		http.Error(w, "missing repo query parameter", http.StatusBadRequest)
		return
	}

	repoPath, err := resolveRepoPath(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/x-diff; charset=utf-8")
	if err := writeRepoDiff(r.Context(), w, repoPath); err != nil {
		log.Printf("Failed to write diff for %s: %v", repoPath, err)
	}
}

func writeRepoDiff(ctx context.Context, w io.Writer, repoPath string) error {
	repoName := repoDisplayName(repoPath)
	return runCommand(ctx, repoPath, w, w, "bash", "-c", gitDiffScript(repoName))
}

func serveDiffsText(w http.ResponseWriter, r *http.Request) {
	writer := &maxSizeWriter{Writer: w, maxSize: 5 * 1024 * 1024}

	repos, err := findGitRepos(".")
	if err != nil {
		http.Error(w, "Failed to find git repositories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/x-diff; charset=utf-8")

	for _, repoPath := range repos {
		if err := writeRepoDiff(r.Context(), writer, repoPath); err != nil {
			continue
		}
	}
}

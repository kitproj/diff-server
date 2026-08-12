package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxRepoDepth = 2

func pathDepth(rel string) int {
	rel = filepath.ToSlash(rel)
	if rel == "." {
		return 0
	}
	return strings.Count(rel, "/") + 1
}

func findGitRepos(root string) ([]string, error) {
	var repos []string
	root = filepath.Clean(root)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repoRel, err := filepath.Rel(root, repoPath)
			if err != nil {
				return nil
			}
			if pathDepth(repoRel) <= maxRepoDepth {
				repos = append(repos, repoPath)
			}
			return filepath.SkipDir
		}

		if info.IsDir() && pathDepth(rel) > maxRepoDepth {
			return filepath.SkipDir
		}

		return nil
	})

	return repos, err
}

func repoDisplayName(repoPath string) string {
	relPath, err := filepath.Rel(".", repoPath)
	if err != nil || relPath == "." {
		return filepath.Base(repoPath)
	}
	return relPath
}

func resolveRepoPath(name string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if name == filepath.Base(cwd) {
		if _, err := os.Stat(filepath.Join(".", ".git")); err == nil {
			return ".", nil
		}
	}

	candidate := filepath.Clean(name)
	if candidate == ".." || strings.HasPrefix(candidate, "../") {
		return "", fmt.Errorf("unknown repo: %s", name)
	}

	if _, err := os.Stat(filepath.Join(candidate, ".git")); err != nil {
		return "", fmt.Errorf("unknown repo: %s", name)
	}

	if pathDepth(candidate) > maxRepoDepth {
		return "", fmt.Errorf("unknown repo: %s", name)
	}

	if repoDisplayName(candidate) != name {
		return "", fmt.Errorf("unknown repo: %s", name)
	}

	return candidate, nil
}

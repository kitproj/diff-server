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
	repos, err := findGitRepos(".")
	if err != nil {
		return "", err
	}
	for _, repoPath := range repos {
		if repoDisplayName(repoPath) == name {
			return repoPath, nil
		}
	}
	return "", fmt.Errorf("unknown repo: %s", name)
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/", handleRoot)
	
	fmt.Printf("Starting server on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	
	repos := findGitRepos(cwd)
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body><h1>Git Diffs</h1>"))
	
	for _, repo := range repos {
		w.Write([]byte(fmt.Sprintf("<h2>%s</h2>", repo)))
		diff := getGitDiff(repo)
		w.Write([]byte("<pre>"))
		w.Write([]byte(diff))
		w.Write([]byte("</pre>"))
	}
	
	w.Write([]byte("</body></html>"))
}

func findGitRepos(root string) []string {
	var repos []string
	
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
	
	return repos
}

func getGitDiff(repoPath string) string {
	cmd := exec.Command("git", "diff")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Sprintf("Error getting diff: %s", err)
	}
	
	if len(output) == 0 {
		cmd = exec.Command("git", "diff", "HEAD")
		cmd.Dir = repoPath
		output, _ = cmd.CombinedOutput()
	}
	
	if len(output) == 0 {
		return "No changes"
	}
	
	return strings.TrimSpace(string(output))
}

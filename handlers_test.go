package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to set git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to set git name: %v", err)
	}
}

func TestServeDiffsHTML(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")

	w := httptest.NewRecorder()
	diffsHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected content-type text/html, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Errorf("expected HTML document")
	}
}

func TestServeDiffsText_PWDIsGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	setupTestGitRepo(t, tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	if err := os.WriteFile(testFile, []byte("modified content\n"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/x-diff")

	w := httptest.NewRecorder()
	diffsHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/x-diff") {
		t.Errorf("expected content-type text/x-diff, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "diff --git") {
		t.Errorf("expected diff output, got: %s", body)
	}
	if !strings.Contains(body, "test.txt") {
		t.Errorf("expected test.txt in diff, got: %s", body)
	}
}

func TestServeDiffsText_GitSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()

	subDir := filepath.Join(tmpDir, "subproject")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	setupTestGitRepo(t, subDir)

	testFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(testFile, []byte("sub content\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd := exec.Command("git", "add", "sub.txt")
	cmd.Dir = subDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "sub commit")
	cmd.Dir = subDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	if err := os.WriteFile(testFile, []byte("modified sub\n"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/x-diff")

	w := httptest.NewRecorder()
	diffsHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "diff --git") {
		t.Errorf("expected diff output, got: %s", body)
	}
	if !strings.Contains(body, "sub.txt") {
		t.Errorf("expected sub.txt in diff, got: %s", body)
	}
	if !strings.Contains(body, "subproject") {
		t.Errorf("expected subproject in diff path, got: %s", body)
	}
}

func TestServeDiffsText_LargeDiffTruncation(t *testing.T) {
	tmpDir := t.TempDir()

	setupTestGitRepo(t, tmpDir)

	testFile := filepath.Join(tmpDir, "large.txt")
	largeContent := strings.Repeat("x", 1024*1024)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd := exec.Command("git", "add", "large.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "large commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	modifiedContent := strings.Repeat("y", 1024*1024)
	if err := os.WriteFile(testFile, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	anotherFile := filepath.Join(tmpDir, "another.txt")
	moreContent := strings.Repeat("z", 5*1024*1024)
	if err := os.WriteFile(anotherFile, []byte(moreContent), 0644); err != nil {
		t.Fatalf("failed to write another file: %v", err)
	}

	cmd = exec.Command("git", "add", "another.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "another commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	yetAnotherContent := strings.Repeat("w", 5*1024*1024)
	if err := os.WriteFile(anotherFile, []byte(yetAnotherContent), 0644); err != nil {
		t.Fatalf("failed to modify another file: %v", err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/x-diff")

	w := httptest.NewRecorder()
	diffsHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if len(body) > 5*1024*1024 {
		t.Errorf("expected diff to be truncated to 5MB, got %d bytes", len(body))
	}

	if len(body) == 0 {
		t.Error("expected some diff output")
	}
}

func TestServeDiffsText_NonGitDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a non-git directory with a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/x-diff")
	
	w := httptest.NewRecorder()
	diffsHandler(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	
	// Should return empty response since there are no git repos
	// This tests that the handler completes successfully even with no repos
}


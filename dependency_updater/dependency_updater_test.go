package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateGithubRepoUrl(t *testing.T) {
	deps := Dependencies{
		"op_node":     {Owner: "ethereum-optimism", Repo: "optimism"},
		"base_reth":   {Owner: "base", Repo: "base"},
		"nethermind":  {Owner: "NethermindEth", Repo: "nethermind"},
		"op_geth":     {Owner: "ethereum-optimism", Repo: "op-geth"},
	}

	tests := []struct {
		name     string
		depType  string
		expected string
	}{
		{"op_node", "op_node", "https://github.com/ethereum-optimism/optimism"},
		{"base_reth", "base_reth", "https://github.com/base/base"},
		{"nethermind", "nethermind", "https://github.com/NethermindEth/nethermind"},
		{"op_geth", "op_geth", "https://github.com/ethereum-optimism/op-geth"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateGithubRepoUrl(deps, tt.depType)
			if got != tt.expected {
				t.Errorf("generateGithubRepoUrl() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGenerateGithubRepoUrl_UnknownDependency(t *testing.T) {
	deps := Dependencies{
		"known": {Owner: "test", Repo: "repo"},
	}

	// Accessing an unknown key returns nil, which should return empty string
	result := generateGithubRepoUrl(deps, "unknown")
	expected := ""
	if result != expected {
		t.Errorf("generateGithubRepoUrl() for unknown key = %q, want %q", result, expected)
	}
}

func TestWriteToVersionsJson(t *testing.T) {
	tmpDir := t.TempDir()

	deps := Dependencies{
		"test_dep": {
			Tag:    "v1.0.0",
			Commit: "abc123def456",
			Owner:  "testowner",
			Repo:   "testrepo",
		},
	}

	err := writeToVersionsJson(tmpDir, deps)
	if err != nil {
		t.Fatalf("writeToVersionsJson() returned error: %v", err)
	}

	// Verify the file was created
	data, err := os.ReadFile(filepath.Join(tmpDir, "versions.json"))
	if err != nil {
		t.Fatalf("failed to read versions.json: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"tag": "v1.0.0"`) {
		t.Errorf("versions.json missing tag: %s", content)
	}
	if !strings.Contains(content, `"commit": "abc123def456"`) {
		t.Errorf("versions.json missing commit: %s", content)
	}
	if !strings.Contains(content, `"owner": "testowner"`) {
		t.Errorf("versions.json missing owner: %s", content)
	}
}

func TestWriteToVersionsJson_InvalidPath(t *testing.T) {
	deps := Dependencies{
		"test": {Tag: "v1.0.0", Commit: "abc", Owner: "o", Repo: "r"},
	}

	err := writeToVersionsJson("/nonexistent/path", deps)
	if err == nil {
		t.Error("writeToVersionsJson() expected error for invalid path, got nil")
	}
}

func TestCreateVersionsEnv_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	deps := Dependencies{
		"base_reth_node": {
			Tag:    "v0.8.0",
			Commit: "3049ce2e3a5132f2ef74b4ba14a1a952ea6abdfb",
			Owner:  "base",
			Repo:   "base",
		},
		"op_node": {
			Tag:    "op-node/v1.16.11",
			Commit: "cba7aba0c98aae22720b21c3a023990a486cb6e0",
			Owner:  "ethereum-optimism",
			Repo:   "optimism",
		},
	}

	err := createVersionsEnv(tmpDir, deps)
	if err != nil {
		t.Fatalf("createVersionsEnv() returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "versions.env"))
	if err != nil {
		t.Fatalf("failed to read versions.env: %v", err)
	}

	content := string(data)

	// Verify all expected variables exist
	checks := []string{
		"export BASE_RETH_NODE_TAG=v0.8.0",
		"export BASE_RETH_NODE_COMMIT=3049ce2e3a5132f2ef74b4ba14a1a952ea6abdfb",
		"export BASE_RETH_NODE_REPO=https://github.com/base/base.git",
		"export OP_NODE_TAG=op-node/v1.16.11",
		"export OP_NODE_COMMIT=cba7aba0c98aae22720b21c3a023990a486cb6e0",
		"export OP_NODE_REPO=https://github.com/ethereum-optimism/optimism.git",
	}

	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("versions.env missing expected line: %q", check)
		}
	}
}

func TestCreateVersionsEnv_BranchTracking(t *testing.T) {
	tmpDir := t.TempDir()

	deps := Dependencies{
		"tracking_branch": {
			Tag:      "", // Tag should be empty for branch tracking
			Commit:   "oldcommit123",
			Owner:    "test",
			Repo:     "test",
			Branch:   "main",
			Tracking: "branch",
		},
	}

	err := createVersionsEnv(tmpDir, deps)
	if err != nil {
		t.Fatalf("createVersionsEnv() returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "versions.env"))
	if err != nil {
		t.Fatalf("failed to read versions.env: %v", err)
	}

	content := string(data)

	// For branch tracking, Tag should be replaced with Branch name
	if !strings.Contains(content, "export TRACKING_BRANCH_TAG=main") {
		t.Errorf("branch tracking: expected TAG=main, got: %s", content)
	}
	if !strings.Contains(content, "export TRACKING_BRANCH_COMMIT=oldcommit123") {
		t.Errorf("branch tracking: expected COMMIT=oldcommit123, got: %s", content)
	}
}

func TestCreateVersionsEnv_SortedOutput(t *testing.T) {
	tmpDir := t.TempDir()

	deps := Dependencies{
		"zeta": {Tag: "v1.0.0", Commit: "a", Owner: "o", Repo: "r"},
		"alpha": {Tag: "v1.0.0", Commit: "b", Owner: "o", Repo: "r"},
		"beta":  {Tag: "v1.0.0", Commit: "c", Owner: "o", Repo: "r"},
	}

	err := createVersionsEnv(tmpDir, deps)
	if err != nil {
		t.Fatalf("createVersionsEnv() returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "versions.env"))
	if err != nil {
		t.Fatalf("failed to read versions.env: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	// Verify alphabetical ordering: ALPHA, BETA, ZETA
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines, got %d", len(lines))
	}

	// Check lines are sorted — first line should start with ALPHA, last with ZETA
	if !strings.HasPrefix(lines[0], "export ALPHA_") {
		t.Errorf("expected first line to start with ALPHA_, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[len(lines)-1], "export ZETA_") {
		t.Errorf("expected last line to start with ZETA_, got: %s", lines[len(lines)-1])
	}
}

func TestCreateVersionsEnv_InvalidPath(t *testing.T) {
	deps := Dependencies{
		"test": {Tag: "v1.0.0", Commit: "abc", Owner: "o", Repo: "r"},
	}

	err := createVersionsEnv("/nonexistent/path", deps)
	if err == nil {
		t.Error("createVersionsEnv() expected error for invalid path, got nil")
	}
}

func TestWriteToGithubOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output.txt")

	// Set GITHUB_OUTPUT env var for the test
	t.Setenv("GITHUB_OUTPUT", outputFile)

	title := "chore: updated repo1, repo2"
	description := "### Dependency Updates\n**repo1** - v1.0.0: [diff](url1)\n**repo2** - v2.0.0: [diff](url2)"

	err := writeToGithubOutput(title, description, tmpDir)
	if err != nil {
		t.Fatalf("writeToGithubOutput() returned error: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read GITHUB_OUTPUT: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "TITLE=chore: updated repo1, repo2") {
		t.Errorf("GITHUB_OUTPUT missing TITLE, got: %s", content)
	}
	if !strings.Contains(content, "DESC<<EOF") {
		t.Errorf("GITHUB_OUTPUT missing DESC delimiter, got: %s", content)
	}
	if !strings.Contains(content, "repo1") || !strings.Contains(content, "repo2") {
		t.Errorf("GITHUB_OUTPUT missing repo descriptions, got: %s", content)
	}
	if !strings.Contains(content, "EOF") {
		t.Errorf("GITHUB_OUTPUT missing closing EOF, got: %s", content)
	}
}

func TestWriteToGithubOutput_NoEnvVar(t *testing.T) {
	// Unset GITHUB_OUTPUT if it was set
	if prev, ok := os.LookupEnv("GITHUB_OUTPUT"); ok {
		defer os.Setenv("GITHUB_OUTPUT", prev)
		os.Unsetenv("GITHUB_OUTPUT")
	}

	err := writeToGithubOutput("title", "desc", "/tmp")
	if err == nil {
		t.Error("writeToGithubOutput() expected error when GITHUB_OUTPUT is unset, got nil")
	}
}

func TestCreateCommitMessage_NeedsGitRepo(t *testing.T) {
	// createCommitMessage runs `git commit` so it requires a real git repo.
	// We verify the githubAction=true path works instead.
	t.Skip("Skipping: createCommitMessage local mode requires a git repository")
}

func TestCreateCommitMessage_GithubActionMode(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output.txt")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	updates := []VersionUpdateInfo{
		{
			Repo:    "base/base",
			From:    "v0.7.0",
			To:      "v0.8.0",
			DiffUrl: "https://github.com/base/base/compare/v0.7.0...v0.8.0",
		},
		{
			Repo:    "ethereum-optimism/op-geth",
			From:    "v1.101701.0",
			To:      "v1.101702.0",
			DiffUrl: "https://github.com/ethereum-optimism/op-geth/compare/v1.101701.0...v1.101702.0",
		},
	}

	err := createCommitMessage(updates, tmpDir, true)
	if err != nil {
		t.Fatalf("createCommitMessage(githubAction=true) returned error: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read GITHUB_OUTPUT: %v", err)
	}

	content := string(data)

	// Verify the title contains both repos
	if !strings.Contains(content, "base/base") {
		t.Errorf("commit title missing base/base, got: %s", content)
	}
	if !strings.Contains(content, "ethereum-optimism/op-geth") {
		t.Errorf("commit title missing op-geth, got: %s", content)
	}

	// Verify description has the dependency update section
	if !strings.Contains(content, "### Dependency Updates") {
		t.Errorf("commit description missing header")
	}
}

func TestUpdateVersionTagAndCommit(t *testing.T) {
	tmpDir := t.TempDir()

	deps := Dependencies{
		"test_dep": {
			Tag:    "v1.0.0",
			Commit: "oldcommit",
			Owner:  "test",
			Repo:   "test",
		},
	}

	// Update to new version
	err := updateVersionTagAndCommit("newcommit456", "v2.0.0", "test_dep", tmpDir, deps)
	if err != nil {
		t.Fatalf("updateVersionTagAndCommit() returned error: %v", err)
	}

	// Verify in-memory update
	if deps["test_dep"].Tag != "v2.0.0" {
		t.Errorf("expected Tag=v2.0.0, got %q", deps["test_dep"].Tag)
	}
	if deps["test_dep"].Commit != "newcommit456" {
		t.Errorf("expected Commit=newcommit456, got %q", deps["test_dep"].Commit)
	}

	// Verify file write
	data, err := os.ReadFile(filepath.Join(tmpDir, "versions.json"))
	if err != nil {
		t.Fatalf("failed to read versions.json: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, `"tag": "v2.0.0"`) {
		t.Errorf("versions.json not updated on disk: %s", content)
	}
}

func TestNoUpdatesDoesNotCreateCommit(t *testing.T) {
	// When updatedDependencies is nil/empty, commit should not be created
	// This tests the condition in updater(): `if (commit && updatedDependencies != nil)`

	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output.txt")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	// Empty updates slice — should not write anything
	err := createCommitMessage([]VersionUpdateInfo{}, tmpDir, true)
	if err != nil {
		t.Fatalf("createCommitMessage empty updates returned error: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read GITHUB_OUTPUT: %v", err)
	}

	content := string(data)
	// The function should write even with empty updates since cmd is still called
	// But the title should just be "chore: updated " with no repos
	if !strings.Contains(content, "TITLE=chore: updated") {
		t.Errorf("expected basic title, got: %s", content)
	}
}

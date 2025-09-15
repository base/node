package main

import (
	"reflect"
	"testing"
)

func TestBuildCommitMessagePartsDeterministic(t *testing.T) {
	t.Helper()

	deps := []VersionUpdateInfo{
		{Repo: "op_geth", To: "v1.0.0", DiffUrl: "diff-geth"},
		{Repo: "node-reth", To: "v0.1.0", DiffUrl: "diff-reth"},
		{Repo: "optimism", To: "op-node/v1.13.4", DiffUrl: "diff-node"},
	}

	title, description := buildCommitMessageParts(deps)

	wantTitle := "chore: updated node-reth, op_geth, optimism"
	if title != wantTitle {
		t.Fatalf("unexpected commit title: got %q want %q", title, wantTitle)
	}

	wantDescription := "Updated dependencies for: node-reth => v0.1.0 (diff-reth) op_geth => v1.0.0 (diff-geth) optimism => op-node/v1.13.4 (diff-node)"
	if description != wantDescription {
		t.Fatalf("unexpected commit description: got %q want %q", description, wantDescription)
	}
}

func TestBuildCommitMessagePartsDoesNotMutateInput(t *testing.T) {
	t.Helper()

	deps := []VersionUpdateInfo{
		{Repo: "op_geth", To: "v1.0.0", DiffUrl: "diff-geth"},
		{Repo: "node-reth", To: "v0.1.0", DiffUrl: "diff-reth"},
	}

	original := make([]VersionUpdateInfo, len(deps))
	copy(original, deps)

	buildCommitMessageParts(deps)

	if !reflect.DeepEqual(deps, original) {
		t.Fatalf("buildCommitMessageParts mutated input slice: got %+v want %+v", deps, original)
	}
}

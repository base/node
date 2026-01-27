package main

import (
	"testing"
)

func TestNormalizeRCFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v0.3.0-rc1", "v0.3.0-rc.1"},
		{"v0.3.0-rc.1", "v0.3.0-rc.1"},
		{"v0.3.0-rc-1", "v0.3.0-rc.1"},
		{"v0.3.0-RC1", "v0.3.0-rc.1"},
		{"v0.3.0-rc12", "v0.3.0-rc.12"},
		{"v0.3.0", "v0.3.0"},
		{"v0.3.0-alpha", "v0.3.0-alpha"},
		{"v0.3.0-beta.1", "v0.3.0-beta.1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeRCFormat(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeRCFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		tag       string
		tagPrefix string
		wantErr   bool
	}{
		// Standard versions
		{"v0.2.2", "", false},
		{"v0.3.0", "", false},
		{"1.35.3", "", false}, // nethermind style - no v prefix

		// RC versions
		{"v0.3.0-rc1", "", false},
		{"v0.3.0-rc.1", "", false},
		{"v0.3.0-rc-1", "", false},
		{"v0.3.0-rc.2", "", false},

		// With tagPrefix
		{"op-node/v1.16.2", "op-node", false},
		{"op-node/v1.16.3-rc1", "op-node", false},

		// Non-standard but parseable
		{"v1.101603.5", "", false}, // op-geth style

		// Invalid
		{"not-a-version", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			_, err := ParseVersion(tt.tag, tt.tagPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion(%q, %q) error = %v, wantErr %v", tt.tag, tt.tagPrefix, err, tt.wantErr)
			}
		})
	}
}

func TestValidateVersionUpgrade(t *testing.T) {
	tests := []struct {
		name       string
		currentTag string
		newTag     string
		tagPrefix  string
		wantErr    bool
	}{
		// Valid upgrades
		{"stable to rc", "v0.2.2", "v0.3.0-rc1", "", false},
		{"rc to rc", "v0.3.0-rc1", "v0.3.0-rc2", "", false},
		{"rc to stable", "v0.3.0-rc2", "v0.3.0", "", false},
		{"stable to stable", "v0.2.2", "v0.3.0", "", false},
		{"patch upgrade", "v0.2.2", "v0.2.3", "", false},
		{"minor upgrade", "v0.2.2", "v0.3.0", "", false},
		{"major upgrade", "v0.2.2", "v1.0.0", "", false},
		{"same version", "v0.2.2", "v0.2.2", "", false},

		// With tagPrefix
		{"prefix upgrade", "op-node/v1.16.2", "op-node/v1.16.3", "op-node", false},
		{"prefix rc upgrade", "op-node/v1.16.2", "op-node/v1.17.0-rc1", "op-node", false},

		// No v prefix (nethermind style)
		{"no v prefix upgrade", "1.35.3", "1.35.4", "", false},

		// Invalid downgrades
		{"downgrade major", "v0.3.0", "v0.2.2", "", true},
		{"downgrade minor", "v0.3.0", "v0.2.9", "", true},
		{"downgrade patch", "v0.3.1", "v0.3.0", "", true},
		{"stable to rc same version", "v0.3.0", "v0.3.0-rc2", "", true},
		{"stable to rc older version", "v0.3.0", "v0.2.0-rc1", "", true},

		// Edge cases
		{"empty current - valid new", "", "v0.3.0", "", false},
		{"empty current - invalid new", "", "not-a-version", "", true},
		{"unparseable current allows update", "not-semver", "v0.3.0", "", false},

		// Unparseable current with unparseable new should fail
		{"unparseable current - unparseable new", "rollup-boost/v0.7.11", "websocket-proxy/v0.0.2", "", true},
		{"unparseable current - valid new", "rollup-boost/v0.7.11", "v0.8.0", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersionUpgrade(tt.currentTag, tt.newTag, tt.tagPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersionUpgrade(%q, %q, %q) error = %v, wantErr %v",
					tt.currentTag, tt.newTag, tt.tagPrefix, err, tt.wantErr)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name      string
		v1        string
		v2        string
		tagPrefix string
		want      int
	}{
		{"v1 less than v2", "v0.2.2", "v0.3.0", "", -1},
		{"v1 greater than v2", "v0.3.0", "v0.2.2", "", 1},
		{"equal versions", "v0.3.0", "v0.3.0", "", 0},
		{"rc less than stable", "v0.3.0-rc1", "v0.3.0", "", -1},
		{"rc1 less than rc2", "v0.3.0-rc1", "v0.3.0-rc2", "", -1},
		{"stable greater than rc", "v0.3.0", "v0.3.0-rc2", "", 1},
		{"with prefix", "op-node/v1.16.2", "op-node/v1.16.3", "op-node", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareVersions(tt.v1, tt.v2, tt.tagPrefix)
			if err != nil {
				t.Errorf("CompareVersions(%q, %q, %q) unexpected error: %v", tt.v1, tt.v2, tt.tagPrefix, err)
				return
			}
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q, %q) = %d, want %d", tt.v1, tt.v2, tt.tagPrefix, got, tt.want)
			}
		})
	}
}

func TestIsReleaseVersion(t *testing.T) {
	tests := []struct {
		tag       string
		tagPrefix string
		want      bool
	}{
		// Stable releases
		{"v1.0.0", "", true},
		{"v0.2.2", "", true},
		{"1.35.3", "", true}, // nethermind style
		{"v1.101603.5", "", true}, // op-geth style

		// With prefix
		{"op-node/v1.16.2", "op-node", true},

		// Pre-release versions (should return false)
		{"v1.0.0-rc1", "", false},
		{"v1.0.0-rc.1", "", false},
		{"v1.0.0-rc-1", "", false},
		{"v1.0.0-synctest.0", "", false},
		{"v1.0.0-alpha", "", false},
		{"v1.0.0-beta.1", "", false},
		{"op-node/v1.16.6-synctest.0", "op-node", false},
		{"op-node/v1.16.3-rc1", "op-node", false},

		// Invalid versions (should return false)
		{"not-a-version", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got := IsReleaseVersion(tt.tag, tt.tagPrefix)
			if got != tt.want {
				t.Errorf("IsReleaseVersion(%q, %q) = %v, want %v", tt.tag, tt.tagPrefix, got, tt.want)
			}
		})
	}
}

func TestIsRCVersion(t *testing.T) {
	tests := []struct {
		tag       string
		tagPrefix string
		want      bool
	}{
		// RC versions
		{"v1.0.0-rc1", "", true},
		{"v1.0.0-rc.1", "", true},
		{"v1.0.0-rc-1", "", true},
		{"v1.0.0-RC1", "", true},
		{"v1.0.0-rc12", "", true},
		{"op-node/v1.16.3-rc1", "op-node", true},
		{"op-node/v1.16.3-rc.2", "op-node", true},

		// Stable releases (not RC)
		{"v1.0.0", "", false},
		{"v0.2.2", "", false},
		{"op-node/v1.16.2", "op-node", false},

		// Other pre-release versions (not RC)
		{"v1.0.0-synctest.0", "", false},
		{"op-node/v1.16.6-synctest.0", "op-node", false},
		{"v1.0.0-alpha", "", false},
		{"v1.0.0-beta.1", "", false},
		{"v1.0.0-alpha.rc1", "", false}, // rc is part of another prerelease

		// Invalid versions
		{"not-a-version", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got := IsRCVersion(tt.tag, tt.tagPrefix)
			if got != tt.want {
				t.Errorf("IsRCVersion(%q, %q) = %v, want %v", tt.tag, tt.tagPrefix, got, tt.want)
			}
		})
	}
}

func TestIsReleaseOrRCVersion(t *testing.T) {
	tests := []struct {
		tag       string
		tagPrefix string
		want      bool
	}{
		// Stable releases - should pass
		{"v1.0.0", "", true},
		{"v0.2.2", "", true},
		{"op-node/v1.16.2", "op-node", true},

		// RC versions - should pass
		{"v1.0.0-rc1", "", true},
		{"v1.0.0-rc.1", "", true},
		{"op-node/v1.16.3-rc1", "op-node", true},

		// Other pre-release versions - should NOT pass
		{"v1.0.0-synctest.0", "", false},
		{"op-node/v1.16.6-synctest.0", "op-node", false},
		{"v1.0.0-alpha", "", false},
		{"v1.0.0-beta.1", "", false},

		// Invalid versions
		{"not-a-version", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got := IsReleaseOrRCVersion(tt.tag, tt.tagPrefix)
			if got != tt.want {
				t.Errorf("IsReleaseOrRCVersion(%q, %q) = %v, want %v", tt.tag, tt.tagPrefix, got, tt.want)
			}
		})
	}
}

func TestRCVersionOrdering(t *testing.T) {
	// Verify that RC versions are ordered correctly
	versions := []string{
		"v0.2.2",
		"v0.3.0-rc.1",
		"v0.3.0-rc.2",
		"v0.3.0",
		"v0.3.1",
	}

	for i := 0; i < len(versions)-1; i++ {
		current := versions[i]
		next := versions[i+1]
		t.Run(current+" -> "+next, func(t *testing.T) {
			err := ValidateVersionUpgrade(current, next, "")
			if err != nil {
				t.Errorf("Expected %s -> %s to be valid upgrade, got error: %v", current, next, err)
			}
		})
	}

	// Verify reverse order is invalid
	for i := len(versions) - 1; i > 0; i-- {
		current := versions[i]
		previous := versions[i-1]
		t.Run(current+" -> "+previous+" (downgrade)", func(t *testing.T) {
			err := ValidateVersionUpgrade(current, previous, "")
			if err == nil {
				t.Errorf("Expected %s -> %s to be invalid downgrade", current, previous)
			}
		})
	}
}

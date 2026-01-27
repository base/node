package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// rcPattern matches various RC formats: -rc1, -rc.1, -rc-1, -RC1, etc.
var rcPattern = regexp.MustCompile(`(?i)-rc[.-]?(\d+)`)

// rcOnlyPattern is used to check if a version contains ONLY an RC prerelease (not -synctest, -alpha, etc.)
var rcOnlyPattern = regexp.MustCompile(`(?i)^-rc[.-]?\d+$`)

// ParseVersion extracts and normalizes a semantic version from a tag string.
// It handles tagPrefix stripping, v-prefix normalization, and RC format normalization.
func ParseVersion(tag string, tagPrefix string) (*semver.Version, error) {
	versionStr := tag

	// Step 1: Strip tagPrefix if present (e.g., "op-node/v1.16.2" -> "v1.16.2")
	if tagPrefix != "" && strings.HasPrefix(tag, tagPrefix) {
		versionStr = strings.TrimPrefix(tag, tagPrefix)
		versionStr = strings.TrimPrefix(versionStr, "/")
	}

	// Step 2: Normalize RC formats to semver-compatible format
	// "-rc1" -> "-rc.1", "-rc-1" -> "-rc.1"
	versionStr = normalizeRCFormat(versionStr)

	// Step 3: Parse using Masterminds/semver (handles v prefix automatically)
	v, err := semver.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version format %q: %w", tag, err)
	}

	return v, nil
}

// normalizeRCFormat converts various RC formats to semver-compatible format.
// Examples: "-rc1" -> "-rc.1", "-rc-2" -> "-rc.2"
func normalizeRCFormat(version string) string {
	return rcPattern.ReplaceAllString(version, "-rc.$1")
}

// ValidateVersionUpgrade checks if transitioning from currentTag to newTag
// is a valid upgrade (not a downgrade).
// Returns nil if valid, error explaining why if invalid.
func ValidateVersionUpgrade(currentTag, newTag, tagPrefix string) error {
	// First-time setup: no current version, any valid version is acceptable
	if currentTag == "" {
		_, err := ParseVersion(newTag, tagPrefix)
		return err
	}

	// Parse current version
	currentVersion, err := ParseVersion(currentTag, tagPrefix)
	if err != nil {
		// Current version unparseable - still validate new version is parseable
		_, newErr := ParseVersion(newTag, tagPrefix)
		return newErr
	}

	// Parse new version
	newVersion, err := ParseVersion(newTag, tagPrefix)
	if err != nil {
		return fmt.Errorf("new version %q is not a valid semver: %w", newTag, err)
	}

	// Check for downgrade
	if newVersion.LessThan(currentVersion) {
		return fmt.Errorf(
			"version downgrade detected: %s -> %s",
			currentTag, newTag,
		)
	}

	return nil
}

// CompareVersions compares two version tags and returns:
// -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
// Returns 0 and error if either version cannot be parsed.
func CompareVersions(v1Tag, v2Tag, tagPrefix string) (int, error) {
	v1, err := ParseVersion(v1Tag, tagPrefix)
	if err != nil {
		return 0, err
	}
	v2, err := ParseVersion(v2Tag, tagPrefix)
	if err != nil {
		return 0, err
	}
	return v1.Compare(v2), nil
}

// IsReleaseVersion returns true if the tag is a stable release (no prerelease suffix).
// Examples:
//   - "v1.0.0" -> true
//   - "v1.0.0-rc1" -> false
//   - "v1.0.0-synctest.0" -> false
func IsReleaseVersion(tag string, tagPrefix string) bool {
	v, err := ParseVersion(tag, tagPrefix)
	if err != nil {
		return false
	}
	return v.Prerelease() == ""
}

// IsRCVersion returns true if the tag is a release candidate version.
// This matches versions with -rc, -rc.N, -rc-N, -rcN suffixes.
// Examples:
//   - "v1.0.0-rc1" -> true
//   - "v1.0.0-rc.2" -> true
//   - "v1.0.0" -> false (stable release, not RC)
//   - "v1.0.0-synctest.0" -> false (not an RC)
//   - "v1.0.0-alpha" -> false (not an RC)
func IsRCVersion(tag string, tagPrefix string) bool {
	v, err := ParseVersion(tag, tagPrefix)
	if err != nil {
		return false
	}
	prerelease := v.Prerelease()
	if prerelease == "" {
		return false
	}
	// Check if the prerelease is ONLY an RC format (e.g., "rc.1", "rc1", "rc-1")
	// We need to check the original format before normalization
	return rcOnlyPattern.MatchString("-" + prerelease)
}

// IsReleaseOrRCVersion returns true if the tag is either a stable release or an RC version.
// This excludes other prereleases like -alpha, -beta, -synctest, etc.
func IsReleaseOrRCVersion(tag string, tagPrefix string) bool {
	return IsReleaseVersion(tag, tagPrefix) || IsRCVersion(tag, tagPrefix)
}

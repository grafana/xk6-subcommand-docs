package docs

import (
	"errors"
	"runtime/debug"
	"strings"
)

// detectK6Version reads build info using the provided function and returns the
// wildcard-mapped version of the go.k6.io/k6 dependency.
func detectK6Version(readBuildInfo func() (*debug.BuildInfo, bool)) (string, error) {
	info, ok := readBuildInfo()
	if !ok {
		return "", errors.New("build info unavailable")
	}

	for _, dep := range info.Deps {
		if dep.Path == "go.k6.io/k6" {
			return MapToWildcard(dep.Version), nil
		}
	}

	return "", errors.New("go.k6.io/k6 dependency not found in build info")
}

// MapToWildcard converts a semver version to a wildcard patch version.
// "v1.5.0" becomes "v1.5.x", "v0.55.2-rc.1" becomes "v0.55.x".
// Pre-release suffixes and build metadata are stripped.
// The "v" prefix is always added if missing, since the k6-docs repository
// uses v-prefixed directory names (e.g. "v1.6.x").
// If the version doesn't contain at least two dots (major.minor.patch),
// it is returned as-is.
func MapToWildcard(version string) string {
	if version == "" {
		return ""
	}

	// Strip pre-release (-...) and build metadata (+...) first.
	// Find the earliest occurrence of either '-' or '+' after the version prefix.
	clean := version
	if idx := strings.IndexAny(clean, "-+"); idx != -1 {
		clean = clean[:idx]
	}

	// Find the last dot to replace patch with "x".
	lastDot := strings.LastIndex(clean, ".")
	if lastDot == -1 {
		return version
	}

	// Ensure there are at least two dots (major.minor.patch).
	prefix := clean[:lastDot]
	if !strings.Contains(prefix, ".") {
		return version
	}

	result := prefix + ".x"

	// Ensure the "v" prefix is present. The k6-docs repo always uses
	// v-prefixed directory names like "v1.6.x".
	if !strings.HasPrefix(result, "v") {
		result = "v" + result
	}

	return result
}

// DetectK6Version is a convenience wrapper that uses the real debug.ReadBuildInfo.
func DetectK6Version() (string, error) {
	return detectK6Version(debug.ReadBuildInfo)
}

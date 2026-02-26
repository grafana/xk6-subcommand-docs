package docs

import (
	"runtime/debug"
	"testing"
)

func TestMapToWildcard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{name: "standard three-part", version: "v1.5.0", want: "v1.5.x"},
		{name: "high patch number", version: "v0.55.2", want: "v0.55.x"},
		{name: "single digit parts", version: "v0.1.0", want: "v0.1.x"},
		{name: "large numbers", version: "v12.345.678", want: "v12.345.x"},
		{name: "pre-release suffix", version: "v1.5.0-rc.1", want: "v1.5.x"},
		{name: "build metadata suffix", version: "v0.55.2+build.123", want: "v0.55.x"},
		{name: "pre-release and build", version: "v1.0.0-alpha+001", want: "v1.0.x"},
		{name: "pseudo-version", version: "v0.0.0-20240101000000-abcdef123456", want: "v0.0.x"},
		{name: "no v prefix", version: "1.5.0", want: "v1.5.x"},
		{name: "only major.minor", version: "v1.5", want: "v1.5"},
		{name: "empty string", version: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MapToWildcard(tt.version)
			if got != tt.want {
				t.Errorf("MapToWildcard(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestDetectK6Version(t *testing.T) {
	t.Parallel()

	t.Run("k6 dependency found", func(t *testing.T) {
		t.Parallel()

		mock := func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Deps: []*debug.Module{
					{Path: "github.com/spf13/cobra", Version: "v1.10.2"},
					{Path: "go.k6.io/k6", Version: "v1.6.0"},
					{Path: "github.com/sirupsen/logrus", Version: "v1.9.3"},
				},
			}, true
		}

		got, err := detectK6Version(mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "v1.6.x" {
			t.Errorf("detectK6Version() = %q, want %q", got, "v1.6.x")
		}
	})

	t.Run("k6 dependency not found", func(t *testing.T) {
		t.Parallel()

		mock := func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Deps: []*debug.Module{
					{Path: "github.com/spf13/cobra", Version: "v1.10.2"},
				},
			}, true
		}

		_, err := detectK6Version(mock)
		if err == nil {
			t.Fatal("expected error when k6 dependency is missing, got nil")
		}
	})

	t.Run("build info unavailable", func(t *testing.T) {
		t.Parallel()

		mock := func() (*debug.BuildInfo, bool) {
			return nil, false
		}

		_, err := detectK6Version(mock)
		if err == nil {
			t.Fatal("expected error when build info is unavailable, got nil")
		}
	})

	t.Run("k6 pre-release version", func(t *testing.T) {
		t.Parallel()

		mock := func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Deps: []*debug.Module{
					{Path: "go.k6.io/k6", Version: "v0.55.2-rc.1"},
				},
			}, true
		}

		got, err := detectK6Version(mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "v0.55.x" {
			t.Errorf("detectK6Version() = %q, want %q", got, "v0.55.x")
		}
	})
}

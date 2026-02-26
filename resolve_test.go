package docs

import "testing"

func TestResolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		// Edge cases
		{name: "nil args", args: nil, want: ""},
		{name: "empty args", args: []string{}, want: ""},

		// Rule 1: input contains "/" → treat as full slug
		{name: "full slug single arg", args: []string{"javascript-api/k6-http/get"}, want: "javascript-api/k6-http/get"},

		// Rule 2: first word matches known category → join with "/"
		{name: "category with sub", args: []string{"using-k6", "scenarios"}, want: "using-k6/scenarios"},
		{name: "category with deep path", args: []string{"using-k6", "k6-options", "reference"}, want: "using-k6/k6-options/reference"},
		{name: "examples category", args: []string{"examples", "websockets"}, want: "examples/websockets"},
		{name: "testing-guides category", args: []string{"testing-guides", "test-types"}, want: "testing-guides/test-types"},
		{name: "single category word", args: []string{"using-k6"}, want: "using-k6"},
		{name: "single category javascript-api", args: []string{"javascript-api"}, want: "javascript-api"},
		{name: "using-k6-browser category", args: []string{"using-k6-browser", "overview"}, want: "using-k6-browser/overview"},
		{name: "results-output category", args: []string{"results-output", "grafana-cloud"}, want: "results-output/grafana-cloud"},
		{name: "reference category", args: []string{"reference", "glossary"}, want: "reference/glossary"},

		// Rule 3: first word not a category → JS API module shortcut
		{name: "http module shortcut", args: []string{"http", "get"}, want: "javascript-api/k6-http/get"},
		{name: "browser module shortcut", args: []string{"browser", "page", "click"}, want: "javascript-api/k6-browser/page/click"},
		{name: "metrics module shortcut", args: []string{"metrics"}, want: "javascript-api/k6-metrics"},
		{name: "ws module shortcut", args: []string{"ws"}, want: "javascript-api/k6-ws"},
		{name: "data module shortcut", args: []string{"data"}, want: "javascript-api/k6-data"},

		// Rule 3 — duplicate "k6-" prefix must be stripped
		{name: "k6-http already prefixed", args: []string{"k6-http", "get"}, want: "javascript-api/k6-http/get"},
		{name: "k6-metrics already prefixed", args: []string{"k6-metrics"}, want: "javascript-api/k6-metrics"},
		{name: "k6-ws already prefixed", args: []string{"k6-ws"}, want: "javascript-api/k6-ws"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Resolve(tt.args)
			if got != tt.want {
				t.Errorf("Resolve(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestResolveWithLookup(t *testing.T) {
	t.Parallel()

	// Simulate a real index: some slugs have k6- prefix, some don't.
	knownSlugs := map[string]bool{
		"javascript-api/k6-http":    true,
		"javascript-api/k6-http/get": true,
		"javascript-api/k6-metrics": true,
		"javascript-api/k6-data":    true,
		"javascript-api/k6-ws":      true,
		"javascript-api/k6-browser": true,
		"javascript-api/k6-browser/page/click": true,
		"javascript-api/k6-crypto":  true, // k6-crypto exists (different from crypto)
		"javascript-api/jslib":      true, // NO k6- prefix
		"javascript-api/crypto":     true, // NO k6- prefix (different from k6-crypto)
		"javascript-api/init-context":  true, // NO k6- prefix
		"javascript-api/error-codes":   true, // NO k6- prefix
		"using-k6/scenarios":        true,
	}
	exists := func(slug string) bool { return knownSlugs[slug] }

	tests := []struct {
		name string
		args []string
		want string
	}{
		// Edge cases still work
		{name: "nil args", args: nil, want: ""},
		{name: "empty args", args: []string{}, want: ""},

		// Rule 1 and 2 unchanged
		{name: "full slug", args: []string{"javascript-api/k6-http/get"}, want: "javascript-api/k6-http/get"},
		{name: "category", args: []string{"using-k6", "scenarios"}, want: "using-k6/scenarios"},

		// Rule 3: k6- prefixed slugs still resolve correctly
		{name: "http shortcut", args: []string{"http", "get"}, want: "javascript-api/k6-http/get"},
		{name: "metrics shortcut", args: []string{"metrics"}, want: "javascript-api/k6-metrics"},
		{name: "data shortcut", args: []string{"data"}, want: "javascript-api/k6-data"},
		{name: "ws shortcut", args: []string{"ws"}, want: "javascript-api/k6-ws"},
		{name: "k6-http already prefixed", args: []string{"k6-http", "get"}, want: "javascript-api/k6-http/get"},

		// Rule 3: unprefixed slugs now resolve correctly (THE BUG FIX)
		{name: "jslib no k6 prefix", args: []string{"jslib"}, want: "javascript-api/jslib"},
		{name: "crypto no k6 prefix", args: []string{"crypto"}, want: "javascript-api/crypto"},
		{name: "init-context no k6 prefix", args: []string{"init-context"}, want: "javascript-api/init-context"},
		{name: "error-codes no k6 prefix", args: []string{"error-codes"}, want: "javascript-api/error-codes"},

		// k6-crypto (with prefix) is different from crypto (without prefix)
		{name: "k6-crypto with prefix", args: []string{"k6-crypto"}, want: "javascript-api/k6-crypto"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ResolveWithLookup(tt.args, exists)
			if got != tt.want {
				t.Errorf("ResolveWithLookup(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

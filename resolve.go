package docs

import "strings"

// categories lists the known documentation category prefixes.
var categories = []string{
	"javascript-api", "using-k6", "using-k6-browser",
	"testing-guides", "examples", "results-output", "reference",
}

// Resolve converts CLI args into a canonical documentation slug.
// It always assumes the k6- prefix for Rule 3 (JS API shortcuts).
// Use [ResolveWithLookup] when an index is available to handle
// slugs that don't carry the k6- prefix (e.g. jslib, crypto).
func Resolve(args []string) string {
	return ResolveWithLookup(args, nil)
}

// ResolveWithLookup converts CLI args into a canonical documentation slug.
// When exists is non-nil, Rule 3 uses it to disambiguate javascript-api
// children that may or may not carry the k6- prefix.
//
// If the user typed "k6-http", it always resolves to javascript-api/k6-http.
// If the user typed a bare name like "jslib" or "crypto", the function tries
// the unprefixed slug first (javascript-api/crypto), then falls back to the
// k6-prefixed form (javascript-api/k6-crypto). This handles pages like jslib,
// crypto, init-context, and error-codes that don't use the k6- prefix.
func ResolveWithLookup(args []string, exists func(string) bool) string {
	if len(args) == 0 {
		return ""
	}

	// Rule 1: if any arg contains "/", treat as a full slug.
	for _, a := range args {
		if strings.Contains(a, "/") {
			return strings.Join(args, "/")
		}
	}

	// Rule 2: first word matches a known category prefix → join all words.
	for _, cat := range categories {
		if args[0] == cat {
			return strings.Join(args, "/")
		}
	}

	// Rule 3: JS API module shortcut.
	hasK6Prefix := strings.HasPrefix(args[0], "k6-")
	name := strings.TrimPrefix(args[0], "k6-")
	rest := args[1:]
	parts := append([]string{name}, rest...)
	prefixed := "javascript-api/k6-" + strings.Join(parts, "/")

	if exists == nil || hasK6Prefix {
		// No lookup available, or user explicitly typed k6- prefix.
		return prefixed
	}

	// User typed a bare name: try unprefixed first (exact match), then k6- prefixed.
	unprefixed := "javascript-api/" + strings.Join(parts, "/")
	if exists(unprefixed) {
		return unprefixed
	}

	if exists(prefixed) {
		return prefixed
	}

	// Neither found — return prefixed form as default (matches original behavior).
	return prefixed
}

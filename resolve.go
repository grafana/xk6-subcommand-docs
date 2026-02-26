package docs

import "strings"

// categories lists the known documentation category prefixes.
var categories = []string{
	"javascript-api", "using-k6", "using-k6-browser",
	"testing-guides", "examples", "results-output", "reference",
}

// Resolve converts CLI args into a canonical documentation slug.
func Resolve(args []string) string {
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

	// Rule 3: JS API module shortcut → prepend "javascript-api/k6-".
	return "javascript-api/k6-" + strings.Join(args, "/")
}

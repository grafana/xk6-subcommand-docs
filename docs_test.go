package docs

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestCache creates a temp directory with sections.json and markdown files,
// returning the cache dir path and a loaded Index.
func setupTestCache(t *testing.T) (string, *Index) {
	t.Helper()

	dir := t.TempDir()

	sections := []Section{
		{
			Slug:        "javascript-api",
			RelPath:     "javascript-api/_index.md",
			Title:       "JavaScript API",
			Description: "k6 JavaScript API reference.",
			Weight:      1,
			Category:    "javascript-api",
			Children:    []string{"javascript-api/k6-http"},
			IsIndex:     true,
		},
		{
			Slug:        "javascript-api/k6-http",
			RelPath:     "javascript-api/k6-http/_index.md",
			Title:       "k6/http",
			Description: "HTTP module for k6.",
			Weight:      1,
			Category:    "javascript-api",
			Children:    []string{"javascript-api/k6-http/get", "javascript-api/k6-http/post"},
			IsIndex:     true,
		},
		{
			Slug:        "javascript-api/k6-http/get",
			RelPath:     "javascript-api/k6-http/get.md",
			Title:       "get",
			Description: "Make an HTTP GET request.",
			Weight:      1,
			Category:    "javascript-api",
			Children:    nil,
			IsIndex:     false,
		},
		{
			Slug:        "javascript-api/k6-http/post",
			RelPath:     "javascript-api/k6-http/post.md",
			Title:       "post",
			Description: "Make an HTTP POST request.",
			Weight:      2,
			Category:    "javascript-api",
			Children:    nil,
			IsIndex:     false,
		},
		{
			Slug:        "using-k6",
			RelPath:     "using-k6/_index.md",
			Title:       "Using k6",
			Description: "Learn how to use k6.",
			Weight:      2,
			Category:    "using-k6",
			Children:    []string{"using-k6/scenarios"},
			IsIndex:     true,
		},
		{
			Slug:        "using-k6/scenarios",
			RelPath:     "using-k6/scenarios.md",
			Title:       "Scenarios",
			Description: "Configure test scenarios.",
			Weight:      1,
			Category:    "using-k6",
			Children:    nil,
			IsIndex:     false,
		},
		{
			Slug:        "examples",
			RelPath:     "examples/_index.md",
			Title:       "Examples",
			Description: "Example k6 scripts.",
			Weight:      3,
			Category:    "examples",
			Children:    []string{"examples/websockets"},
			IsIndex:     true,
		},
		{
			Slug:        "examples/websockets",
			RelPath:     "examples/websockets.md",
			Title:       "WebSockets",
			Description: "WebSocket load testing examples.",
			Weight:      1,
			Category:    "examples",
			Children:    nil,
			IsIndex:     false,
		},
	}

	idx := &Index{
		Version:  "v0.55.x",
		Sections: sections,
	}

	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("marshal index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sections.json"), data, 0o644); err != nil {
		t.Fatalf("write sections.json: %v", err)
	}

	// Create markdown files. Content includes H1 headings after frontmatter
	// to match how the prepare command generates cached content.
	mdFiles := map[string]string{
		"javascript-api/_index.md":          "---\ntitle: JavaScript API\n---\n# JavaScript API\n\nThe JavaScript API reference.\n",
		"javascript-api/k6-http/_index.md":  "---\ntitle: k6/http\n---\n# k6/http\n\nThe HTTP module.\n",
		"javascript-api/k6-http/get.md":     "---\ntitle: get\n---\n## http.get(url)\n\nMake a GET request.\n",
		"javascript-api/k6-http/post.md":    "---\ntitle: post\n---\n## http.post(url, body)\n\nMake a POST request.\n",
		"using-k6/_index.md":                "---\ntitle: Using k6\n---\n# Using k6\n\nGuide to using k6.\n",
		"using-k6/scenarios.md":             "---\ntitle: Scenarios\n---\n# Scenarios\n\nScenarios let you configure execution.\n",
		"examples/_index.md":                "---\ntitle: Examples\n---\n# Examples\n\nExample scripts.\n",
		"examples/websockets.md":            "---\ntitle: WebSockets\n---\n# WebSockets\n\nWebSocket example content.\n",
	}

	for relPath, content := range mdFiles {
		fullPath := filepath.Join(dir, "markdown", relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(fullPath), err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", fullPath, err)
		}
	}

	// best_practices.md lives at the cache root (same as cmd/prepare output).
	bpPath := filepath.Join(dir, "best_practices.md")
	if err := os.WriteFile(bpPath, []byte("---\ntitle: Best Practices\n---\nFollow these best practices for k6.\n"), 0o644); err != nil {
		t.Fatalf("write best_practices.md: %v", err)
	}

	// Reload from disk so bySlug map is built.
	idx, err = LoadIndex(dir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}

	return dir, idx
}

func TestChildName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		childSlug  string
		parentSlug string
		want       string
	}{
		{"javascript-api/k6-http", "javascript-api", "k6-http"},
		{"javascript-api/k6-http/get", "javascript-api/k6-http", "get"},
		{"javascript-api/k6-http/post", "javascript-api/k6-http", "post"},
		{"using-k6/scenarios", "using-k6", "scenarios"},
		{"examples/websockets", "examples", "websockets"},
		// When parent doesn't match as prefix, fall back to last segment.
		{"javascript-api/k6-http/get", "", "get"},
		{"using-k6/scenarios", "other", "scenarios"},
		// No slash at all returns the slug itself.
		{"toplevel", "", "toplevel"},
		{"toplevel", "other", "toplevel"},
	}

	for _, tt := range tests {
		t.Run(tt.childSlug+"_parent_"+tt.parentSlug, func(t *testing.T) {
			t.Parallel()
			got := childName(tt.childSlug, tt.parentSlug)
			if got != tt.want {
				t.Errorf("childName(%q, %q) = %q, want %q", tt.childSlug, tt.parentSlug, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		max   int
		want  string
	}{
		{"short", 80, "short"},
		{"exactly ten", 11, "exactly ten"},
		{"this string is longer than ten", 10, "this st..."},
		{"", 80, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := truncate(tt.input, tt.max)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
			}
		})
	}
}

func TestPrintTOC(t *testing.T) {
	_, idx := setupTestCache(t)

	var buf bytes.Buffer
	printTOC(&buf, idx, "v0.55.x")
	out := buf.String()

	// Check header.
	if !strings.Contains(out, "k6 Documentation (v0.55.x)") {
		t.Error("printTOC: missing header with version")
	}
	if !strings.Contains(out, "Use: k6 x docs <topic>") {
		t.Error("printTOC: missing usage line")
	}

	// Check categories appear.
	if !strings.Contains(out, "## JavaScript API") {
		t.Error("printTOC: missing JavaScript API category")
	}
	if !strings.Contains(out, "## Using k6") {
		t.Error("printTOC: missing Using k6 category")
	}
	if !strings.Contains(out, "## Examples") {
		t.Error("printTOC: missing Examples category")
	}

	// Check children listed under JavaScript API use child names (relative to parent).
	if !strings.Contains(out, "k6-http") {
		t.Error("printTOC: missing 'k6-http' child name under JavaScript API")
	}
}

func TestPrintSection(t *testing.T) {
	cacheDir, idx := setupTestCache(t)

	t.Run("section with content and children", func(t *testing.T) {
		sec, ok := idx.Lookup("javascript-api/k6-http")
		if !ok {
			t.Fatal("Lookup(javascript-api/k6-http): not found")
		}

		var buf bytes.Buffer
		printSection(&buf, idx, sec, cacheDir, "v0.55.x")
		out := buf.String()

		// Content from markdown file (frontmatter stripped by Transform).
		if !strings.Contains(out, "The HTTP module.") {
			t.Error("printSection: missing markdown content")
		}

		// Subtopics footer.
		if !strings.Contains(out, "---") {
			t.Error("printSection: missing separator")
		}
		if !strings.Contains(out, "Subtopics:") {
			t.Error("printSection: missing Subtopics line")
		}
		if !strings.Contains(out, "get") {
			t.Error("printSection: missing 'get' in subtopics")
		}
		if !strings.Contains(out, "post") {
			t.Error("printSection: missing 'post' in subtopics")
		}
		if !strings.Contains(out, "Use: k6 x docs k6-http <subtopic>") {
			t.Error("printSection: missing usage hint")
		}
	})

	t.Run("leaf section without children", func(t *testing.T) {
		sec, ok := idx.Lookup("javascript-api/k6-http/get")
		if !ok {
			t.Fatal("Lookup(javascript-api/k6-http/get): not found")
		}

		var buf bytes.Buffer
		printSection(&buf, idx, sec, cacheDir, "v0.55.x")
		out := buf.String()

		if !strings.Contains(out, "http.get(url)") {
			t.Error("printSection leaf: missing content")
		}
		if strings.Contains(out, "Subtopics:") {
			t.Error("printSection leaf: should not have subtopics footer")
		}
	})
}

func TestPrintList(t *testing.T) {
	_, idx := setupTestCache(t)

	t.Run("section with children", func(t *testing.T) {
		var buf bytes.Buffer
		printList(&buf, idx, "javascript-api/k6-http")
		out := buf.String()

		if !strings.Contains(out, "k6/http") {
			t.Error("printList: missing section title")
		}
		if !strings.Contains(out, "HTTP module for k6.") {
			t.Error("printList: missing section description")
		}
		if !strings.Contains(out, "get") {
			t.Error("printList: missing child 'get'")
		}
		if !strings.Contains(out, "post") {
			t.Error("printList: missing child 'post'")
		}
	})

	t.Run("section without children", func(t *testing.T) {
		var buf bytes.Buffer
		printList(&buf, idx, "javascript-api/k6-http/get")
		out := buf.String()

		if !strings.Contains(out, "get") {
			t.Error("printList no children: missing title")
		}
		if !strings.Contains(out, "(no subtopics)") {
			t.Error("printList no children: missing 'no subtopics' message")
		}
	})

	t.Run("nonexistent slug", func(t *testing.T) {
		var buf bytes.Buffer
		printList(&buf, idx, "does-not-exist")
		out := buf.String()

		if !strings.Contains(out, "Topic not found") {
			t.Error("printList missing: expected 'Topic not found' message")
		}
	})
}

func TestPrintSearch(t *testing.T) {
	cacheDir, idx := setupTestCache(t)

	t.Run("match in title groups by parent", func(t *testing.T) {
		var buf bytes.Buffer
		printSearch(&buf, idx, "Scenarios", cacheDir)
		out := buf.String()

		if !strings.Contains(out, `Results for "Scenarios"`) {
			t.Error("printSearch: missing results header")
		}
		// Should be grouped under using-k6.
		if !strings.Contains(out, "using-k6:") {
			t.Error("printSearch: missing 'using-k6:' group header")
		}
		if !strings.Contains(out, "scenarios") {
			t.Error("printSearch: missing 'scenarios' result")
		}
	})

	t.Run("match in description shows group with description", func(t *testing.T) {
		var buf bytes.Buffer
		printSearch(&buf, idx, "GET request", cacheDir)
		out := buf.String()

		// Should be grouped under k6-http.
		if !strings.Contains(out, "k6-http:") {
			t.Error("printSearch: missing 'k6-http:' group header")
		}
		if !strings.Contains(out, "get") {
			t.Error("printSearch: missing 'get' result")
		}
	})

	t.Run("match in body content", func(t *testing.T) {
		var buf bytes.Buffer
		printSearch(&buf, idx, "WebSocket example content", cacheDir)
		out := buf.String()

		if !strings.Contains(out, "examples:") {
			t.Error("printSearch: missing 'examples:' group header")
		}
		if !strings.Contains(out, "websockets") {
			t.Error("printSearch: missing 'websockets' result")
		}
	})

	t.Run("groups sorted alphabetically", func(t *testing.T) {
		var buf bytes.Buffer
		// Search for "k6" which should match multiple groups.
		printSearch(&buf, idx, "k6", cacheDir)
		out := buf.String()

		// Verify groups appear in alphabetical order.
		lines := strings.Split(out, "\n")
		var groupHeaders []string
		for _, line := range lines {
			if strings.HasSuffix(line, ":") || (strings.Contains(line, ":") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "Results")) {
				groupHeaders = append(groupHeaders, line)
			}
		}
		for i := 1; i < len(groupHeaders); i++ {
			if groupHeaders[i] < groupHeaders[i-1] {
				t.Errorf("printSearch: groups not sorted: %q before %q", groupHeaders[i-1], groupHeaders[i])
			}
		}
	})

	t.Run("no results", func(t *testing.T) {
		var buf bytes.Buffer
		printSearch(&buf, idx, "zzzznotfound", cacheDir)
		out := buf.String()

		if !strings.Contains(out, "(no results)") {
			t.Error("printSearch no results: missing 'no results' message")
		}
	})
}

func TestPrintBestPractices(t *testing.T) {
	cacheDir, _ := setupTestCache(t)

	var buf bytes.Buffer
	err := printBestPractices(&buf, cacheDir)
	if err != nil {
		t.Fatalf("printBestPractices: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Follow these best practices for k6.") {
		t.Error("printBestPractices: missing content")
	}
	// Frontmatter should be stripped.
	if strings.Contains(out, "---") {
		t.Error("printBestPractices: frontmatter not stripped")
	}
}

func TestPrintBestPracticesMissing(t *testing.T) {
	dir := t.TempDir()
	// No best_practices.md — should return error.
	var buf bytes.Buffer
	err := printBestPractices(&buf, dir)
	if err == nil {
		t.Fatal("printBestPractices: expected error for missing file, got nil")
	}
}

func TestPrintAll(t *testing.T) {
	cacheDir, idx := setupTestCache(t)

	var buf bytes.Buffer
	printAll(&buf, idx, cacheDir, "v0.55.x")
	out := buf.String()

	if !strings.Contains(out, "k6 Documentation (v0.55.x)") {
		t.Error("printAll: missing header")
	}

	// Check that sections are included.
	if !strings.Contains(out, "# JavaScript API") {
		t.Error("printAll: missing JavaScript API section heading")
	}
	if !strings.Contains(out, "# Scenarios") {
		t.Error("printAll: missing Scenarios section heading")
	}
	if !strings.Contains(out, "The HTTP module.") {
		t.Error("printAll: missing HTTP module content")
	}

	// Verify no duplicate headings: each H1 should appear exactly once.
	if strings.Count(out, "# JavaScript API") != 1 {
		t.Errorf("printAll: duplicate '# JavaScript API' heading (count=%d)", strings.Count(out, "# JavaScript API"))
	}
	if strings.Count(out, "# Scenarios") != 1 {
		t.Errorf("printAll: duplicate '# Scenarios' heading (count=%d)", strings.Count(out, "# Scenarios"))
	}
}

func TestCommandIntegration(t *testing.T) {
	cacheDir, _ := setupTestCache(t)

	t.Run("no args shows TOC", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd.Execute: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "k6 Documentation (v0.55.x)") {
			t.Error("integration TOC: missing header")
		}
	})

	t.Run("topic arg shows section", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "using-k6", "scenarios"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd.Execute: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "Scenarios let you configure execution.") {
			t.Error("integration section: missing content")
		}
	})

	t.Run("--list flag shows compact listing", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "--list", "javascript-api/k6-http"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd.Execute: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "get") {
			t.Error("integration --list: missing 'get'")
		}
	})

	t.Run("--all flag prints everything", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "--all"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd.Execute: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "k6 Documentation (v0.55.x)") {
			t.Error("integration --all: missing header")
		}
		if !strings.Contains(out, "# JavaScript API") {
			t.Error("integration --all: missing section")
		}
	})

	t.Run("search subcommand", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "search", "HTTP"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd.Execute: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, `Results for "HTTP"`) {
			t.Error("integration search: missing results header")
		}
	})

	t.Run("best-practices arg", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "best-practices"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd.Execute: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "Follow these best practices for k6.") {
			t.Error("integration best-practices: missing content")
		}
	})

	t.Run("unknown topic returns error", func(t *testing.T) {
		cmd := newCmd(nil)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "nonexistent-topic-xyz"})

		err := cmd.Execute()
		if err == nil {
			t.Fatal("integration unknown topic: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "topic not found") {
			t.Errorf("integration unknown topic: error = %q, want to contain 'topic not found'", err.Error())
		}
	})
}

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	docs "github.com/grafana/xk6-subcommand-docs"
)

// setupMockDocs creates a minimal mock k6-docs directory structure for testing.
// It returns the path that should be used as k6DocsPath.
func setupMockDocs(t *testing.T, version string) string {
	t.Helper()

	root := t.TempDir()
	versionRoot := filepath.Join(root, "docs", "sources", "k6", version)

	// Create version root _index.md (should be skipped).
	writeFile(t, filepath.Join(versionRoot, "_index.md"), `---
title: 'k6 Documentation'
description: 'The k6 documentation.'
weight: 1
---

# k6 Documentation
`)

	// Create shared content.
	writeFile(t, filepath.Join(versionRoot, "shared", "index.md"), `---
headless: true
---
`)
	writeFile(t, filepath.Join(versionRoot, "shared", "javascript-api", "k6-http.md"), `---
title: 'k6/http shared content'
---

The k6/http module contains functionality for performing HTTP transactions.

| Method | Description |
|--------|-------------|
| get    | Issue a GET request. |
`)

	// Create javascript-api category.
	writeFile(t, filepath.Join(versionRoot, "javascript-api", "_index.md"), `---
title: 'JavaScript API'
description: 'The list of k6 modules natively supported in k6 scripts.'
weight: 03
---

# JavaScript API

The list of k6 modules.

{{< section >}}
`)

	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "_index.md"), `---
title: 'k6/http'
description: 'The k6/http module contains functionality for performing HTTP transactions.'
weight: 09
---

# k6/http

{{< docs/shared source="k6" lookup="javascript-api/k6-http.md" version="<K6_VERSION>" >}}
`)

	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "get.md"), `---
title: 'get( url, [params] )'
description: 'Issue an HTTP GET request.'
weight: 10
---

# get( url, [params] )

Make a GET request.

{{< code >}}

`+"```javascript"+`
import http from 'k6/http';

export default function () {
  const res = http.get('https://test.k6.io');
}
`+"```"+`

{{< /code >}}
`)

	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "post.md"), `---
title: 'post( url, [body], [params] )'
description: 'Issue an HTTP POST request.'
weight: 20
---

# post( url, [body], [params] )

Make a POST request.
`)

	// Create using-k6 category.
	writeFile(t, filepath.Join(versionRoot, "using-k6", "_index.md"), `---
title: 'Using k6'
description: 'The using k6 section.'
weight: 05
---

# Using k6
`)

	writeFile(t, filepath.Join(versionRoot, "using-k6", "checks.md"), `---
title: 'Checks'
description: 'Checks validate boolean conditions.'
weight: 400
---

# Checks

Checks validate boolean conditions in your test.

{{< admonition type="note" >}}

When a check fails, the script continues.

{{< /admonition >}}
`)

	writeFile(t, filepath.Join(versionRoot, "using-k6", "thresholds.md"), `---
title: 'Thresholds'
description: 'Thresholds are pass/fail criteria.'
weight: 500
---

# Thresholds

Thresholds are pass/fail criteria for your test metrics.
`)

	// Create reference/glossary (should be included).
	writeFile(t, filepath.Join(versionRoot, "reference", "_index.md"), `---
title: 'Reference'
description: 'k6 reference documentation.'
weight: 100
---

# Reference
`)

	writeFile(t, filepath.Join(versionRoot, "reference", "glossary.md"), `---
title: 'Glossary'
description: 'Technical terms used in k6.'
weight: 07
---

# Glossary

What we talk about when we talk about k6.
`)

	// Create reference/archive.md (should be EXCLUDED — only glossary is included).
	writeFile(t, filepath.Join(versionRoot, "reference", "archive.md"), `---
title: 'Archive'
description: 'k6 archive command.'
weight: 10
---

# Archive
`)

	// Create excluded categories.
	writeFile(t, filepath.Join(versionRoot, "get-started", "_index.md"), `---
title: 'Get Started'
description: 'Getting started with k6.'
weight: 01
---

# Get Started
`)

	writeFile(t, filepath.Join(versionRoot, "extensions", "_index.md"), `---
title: 'Extensions'
description: 'k6 extensions.'
weight: 50
---

# Extensions
`)

	return root
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestRunWithMockDocs(t *testing.T) {
	t.Parallel()

	version := "v0.99.x"
	docsPath := setupMockDocs(t, version)
	outputDir := filepath.Join(t.TempDir(), "output")

	if err := run(version, docsPath, outputDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	// Load and validate sections.json.
	data, err := os.ReadFile(filepath.Join(outputDir, "sections.json"))
	if err != nil {
		t.Fatalf("read sections.json: %v", err)
	}

	var idx docs.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parse sections.json: %v", err)
	}

	if idx.Version != version {
		t.Errorf("Version = %q, want %q", idx.Version, version)
	}

	// Build a slug map for easier assertions.
	bySlug := make(map[string]docs.Section, len(idx.Sections))
	for _, s := range idx.Sections {
		bySlug[s.Slug] = s
	}

	// Check included sections exist.
	expectedSlugs := []string{
		"javascript-api",
		"javascript-api/k6-http",
		"javascript-api/k6-http/get",
		"javascript-api/k6-http/post",
		"using-k6",
		"using-k6/checks",
		"using-k6/thresholds",
		"reference/glossary",
	}
	for _, slug := range expectedSlugs {
		if _, ok := bySlug[slug]; !ok {
			t.Errorf("expected section with slug %q, not found", slug)
		}
	}

	// Check excluded sections are absent.
	excludedSlugs := []string{
		"get-started",
		"extensions",
		"reference",         // _index.md should be excluded (only glossary included)
		"reference/archive", // only glossary from reference
	}
	for _, slug := range excludedSlugs {
		if _, ok := bySlug[slug]; ok {
			t.Errorf("section %q should be excluded", slug)
		}
	}

	// Verify metadata.
	t.Run("javascript-api metadata", func(t *testing.T) {
		s := bySlug["javascript-api"]
		if s.Title != "JavaScript API" {
			t.Errorf("Title = %q, want %q", s.Title, "JavaScript API")
		}
		if s.Category != "javascript-api" {
			t.Errorf("Category = %q, want %q", s.Category, "javascript-api")
		}
		if !s.IsIndex {
			t.Error("IsIndex should be true")
		}
		if s.Weight != 3 {
			t.Errorf("Weight = %d, want 3", s.Weight)
		}
	})

	t.Run("get metadata", func(t *testing.T) {
		s := bySlug["javascript-api/k6-http/get"]
		if s.Title != "get( url, [params] )" {
			t.Errorf("Title = %q, want %q", s.Title, "get( url, [params] )")
		}
		if s.Category != "javascript-api" {
			t.Errorf("Category = %q, want %q", s.Category, "javascript-api")
		}
		if s.IsIndex {
			t.Error("IsIndex should be false")
		}
		if s.Weight != 10 {
			t.Errorf("Weight = %d, want 10", s.Weight)
		}
	})

	// Verify children population.
	t.Run("k6-http children", func(t *testing.T) {
		s := bySlug["javascript-api/k6-http"]
		if len(s.Children) != 2 {
			t.Fatalf("Children count = %d, want 2", len(s.Children))
		}
		// get (weight 10) should come before post (weight 20).
		if s.Children[0] != "javascript-api/k6-http/get" {
			t.Errorf("Children[0] = %q, want %q", s.Children[0], "javascript-api/k6-http/get")
		}
		if s.Children[1] != "javascript-api/k6-http/post" {
			t.Errorf("Children[1] = %q, want %q", s.Children[1], "javascript-api/k6-http/post")
		}
	})

	t.Run("javascript-api children", func(t *testing.T) {
		s := bySlug["javascript-api"]
		if len(s.Children) != 1 {
			t.Fatalf("Children count = %d, want 1 (only k6-http)", len(s.Children))
		}
		if s.Children[0] != "javascript-api/k6-http" {
			t.Errorf("Children[0] = %q, want %q", s.Children[0], "javascript-api/k6-http")
		}
	})

	t.Run("using-k6 children", func(t *testing.T) {
		s := bySlug["using-k6"]
		if len(s.Children) != 2 {
			t.Fatalf("Children count = %d, want 2", len(s.Children))
		}
		// checks (weight 400) before thresholds (weight 500).
		if s.Children[0] != "using-k6/checks" {
			t.Errorf("Children[0] = %q, want %q", s.Children[0], "using-k6/checks")
		}
		if s.Children[1] != "using-k6/thresholds" {
			t.Errorf("Children[1] = %q, want %q", s.Children[1], "using-k6/thresholds")
		}
	})

	t.Run("leaf node has empty children", func(t *testing.T) {
		s := bySlug["using-k6/checks"]
		if s.Children == nil {
			t.Error("Children should be non-nil empty slice")
		}
		if len(s.Children) != 0 {
			t.Errorf("Children count = %d, want 0", len(s.Children))
		}
	})
}

func TestTransformedMarkdownContent(t *testing.T) {
	t.Parallel()

	version := "v0.99.x"
	docsPath := setupMockDocs(t, version)
	outputDir := filepath.Join(t.TempDir(), "output")

	if err := run(version, docsPath, outputDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	t.Run("shared content resolved", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "markdown", "javascript-api", "k6-http", "_index.md"))
		if err != nil {
			t.Fatalf("read k6-http _index.md: %v", err)
		}
		content := string(data)

		if strings.Contains(content, "docs/shared") {
			t.Error("shared shortcode should be resolved")
		}
		if !strings.Contains(content, "k6/http module contains functionality") {
			t.Error("shared content should be inlined")
		}
	})

	t.Run("code tags stripped", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "markdown", "javascript-api", "k6-http", "get.md"))
		if err != nil {
			t.Fatalf("read get.md: %v", err)
		}
		content := string(data)

		if strings.Contains(content, "{{< code >}}") {
			t.Error("code shortcodes should be stripped")
		}
		if !strings.Contains(content, "import http from 'k6/http'") {
			t.Error("code block content should be preserved")
		}
	})

	t.Run("admonition converted", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "markdown", "using-k6", "checks.md"))
		if err != nil {
			t.Fatalf("read checks.md: %v", err)
		}
		content := string(data)

		if strings.Contains(content, "{{< admonition") {
			t.Error("admonition shortcode should be converted")
		}
		if !strings.Contains(content, "> **Note:**") {
			t.Error("admonition should be a blockquote")
		}
	})

	t.Run("frontmatter stripped", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "markdown", "using-k6", "thresholds.md"))
		if err != nil {
			t.Fatalf("read thresholds.md: %v", err)
		}
		content := string(data)

		if strings.Contains(content, "title: 'Thresholds'") {
			t.Error("frontmatter should be stripped from transformed output")
		}
		if !strings.Contains(content, "# Thresholds") {
			t.Error("markdown heading should be preserved")
		}
	})

	t.Run("version placeholder replaced", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "markdown", "javascript-api", "k6-http", "_index.md"))
		if err != nil {
			t.Fatalf("read k6-http _index.md: %v", err)
		}
		content := string(data)

		if strings.Contains(content, "<K6_VERSION>") {
			t.Error("version placeholder should be replaced")
		}
	})
}

func TestBestPracticesWritten(t *testing.T) {
	t.Parallel()

	version := "v0.99.x"
	docsPath := setupMockDocs(t, version)
	outputDir := filepath.Join(t.TempDir(), "output")

	if err := run(version, docsPath, outputDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "best_practices.md"))
	if err != nil {
		t.Fatalf("read best_practices.md: %v", err)
	}

	content := string(data)
	requiredTopics := []string{
		"Test Structure",
		"Performance",
		"Error Handling",
		"Data Management",
		"Authentication",
		"Monitoring",
		"Scenarios",
		"Modules",
		"Browser",
	}
	for _, topic := range requiredTopics {
		if !strings.Contains(content, topic) {
			t.Errorf("best_practices.md should contain %q", topic)
		}
	}

	// Should contain code examples.
	if !strings.Contains(content, "```javascript") {
		t.Error("best_practices.md should contain JavaScript code examples")
	}
}

func TestSlugDerivation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		relPath string
		want    string
	}{
		{"javascript-api/k6-http/_index.md", "javascript-api/k6-http"},
		{"javascript-api/k6-http/get.md", "javascript-api/k6-http/get"},
		{"using-k6/scenarios/_index.md", "using-k6/scenarios"},
		{"using-k6/checks.md", "using-k6/checks"},
		{"examples/_index.md", "examples"},
		{"reference/glossary.md", "reference/glossary"},
	}

	for _, tt := range tests {
		t.Run(tt.relPath, func(t *testing.T) {
			t.Parallel()
			got := slugFromRelPath(tt.relPath)
			if got != tt.want {
				t.Errorf("slugFromRelPath(%q) = %q, want %q", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestCategoryDerivation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		slug string
		want string
	}{
		{"javascript-api/k6-http/get", "javascript-api"},
		{"using-k6/checks", "using-k6"},
		{"examples", "examples"},
		{"reference/glossary", "reference"},
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			t.Parallel()
			got := categoryFromSlug(tt.slug)
			if got != tt.want {
				t.Errorf("categoryFromSlug(%q) = %q, want %q", tt.slug, got, tt.want)
			}
		})
	}
}

func TestIsIncluded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		relPath string
		want    bool
	}{
		{"javascript-api/_index.md", true},
		{"javascript-api/k6-http/get.md", true},
		{"using-k6/checks.md", true},
		{"using-k6-browser/_index.md", true},
		{"testing-guides/api-testing.md", true},
		{"examples/basic.md", true},
		{"results-output/json.md", true},
		{"reference/glossary.md", true},
		{"reference/glossary/_index.md", true},
		{"reference/archive.md", false},
		{"reference/_index.md", false},
		{"reference/k6-rest-api.md", false},
		{"get-started/_index.md", false},
		{"set-up/install.md", false},
		{"extensions/_index.md", false},
		{"grafana-cloud-k6/setup.md", false},
		{"release-notes/v1.5.md", false},
		{"k6-studio/overview.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.relPath, func(t *testing.T) {
			t.Parallel()
			got := isIncluded(tt.relPath)
			if got != tt.want {
				t.Errorf("isIncluded(%q) = %v, want %v", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestParseFrontmatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    frontmatter
	}{
		{
			name:    "typical frontmatter",
			content: "---\ntitle: 'Checks'\ndescription: 'Validate conditions.'\nweight: 400\n---\n\n# Checks",
			want:    frontmatter{Title: "Checks", Description: "Validate conditions.", Weight: 400},
		},
		{
			name:    "no frontmatter",
			content: "# Just markdown",
			want:    frontmatter{},
		},
		{
			name:    "frontmatter with extra fields",
			content: "---\ntitle: 'Test'\naliases:\n  - /old/path\nweight: 5\n---\ncontent",
			want:    frontmatter{Title: "Test", Weight: 5},
		},
		{
			name:    "empty content",
			content: "",
			want:    frontmatter{},
		},
		{
			name:    "duplicate keys keeps first",
			content: "---\ntitle: 'First'\ndescription: 'First desc'\ndescription: 'Second desc'\nweight: 10\n---\n\n# Body",
			want:    frontmatter{Title: "First", Description: "First desc", Weight: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseFrontmatter(tt.content)
			if err != nil {
				t.Fatalf("parseFrontmatter: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDeduplicateYAMLKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no duplicates",
			input: "title: 'Hello'\ndescription: 'World'\nweight: 1",
			want:  "title: 'Hello'\ndescription: 'World'\nweight: 1",
		},
		{
			name:  "duplicate description",
			input: "title: 'Hello'\ndescription: 'First'\ndescription: 'Second'\nweight: 1",
			want:  "title: 'Hello'\ndescription: 'First'\nweight: 1",
		},
		{
			name:  "preserves indented lines",
			input: "title: 'Hello'\naliases:\n  - /old\n  - /older\ntitle: 'Duplicate'",
			want:  "title: 'Hello'\naliases:\n  - /old\n  - /older",
		},
		{
			name:  "preserves comments",
			input: "# comment\ntitle: 'Hello'\ntitle: 'Dup'",
			want:  "# comment\ntitle: 'Hello'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := deduplicateYAMLKeys(tt.input)
			if got != tt.want {
				t.Errorf("deduplicateYAMLKeys():\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestPopulateChildren(t *testing.T) {
	t.Parallel()

	sections := []docs.Section{
		{Slug: "using-k6", IsIndex: true, Weight: 1},
		{Slug: "using-k6/checks", Weight: 400},
		{Slug: "using-k6/thresholds", Weight: 200},
		{Slug: "using-k6/scenarios", IsIndex: true, Weight: 300},
		{Slug: "using-k6/scenarios/executors", IsIndex: true, Weight: 1},
		{Slug: "using-k6/scenarios/executors/shared-iterations", Weight: 1},
	}

	populateChildren(sections)

	// using-k6 should have checks, thresholds, scenarios as direct children.
	// Sorted by weight: thresholds (200), scenarios (300), checks (400).
	parent := sections[0]
	if len(parent.Children) != 3 {
		t.Fatalf("using-k6 children: got %d, want 3", len(parent.Children))
	}
	if parent.Children[0] != "using-k6/thresholds" {
		t.Errorf("Children[0] = %q, want %q", parent.Children[0], "using-k6/thresholds")
	}
	if parent.Children[1] != "using-k6/scenarios" {
		t.Errorf("Children[1] = %q, want %q", parent.Children[1], "using-k6/scenarios")
	}
	if parent.Children[2] != "using-k6/checks" {
		t.Errorf("Children[2] = %q, want %q", parent.Children[2], "using-k6/checks")
	}

	// using-k6/scenarios should have executors as only direct child.
	scenarios := sections[3]
	if len(scenarios.Children) != 1 {
		t.Fatalf("scenarios children: got %d, want 1", len(scenarios.Children))
	}
	if scenarios.Children[0] != "using-k6/scenarios/executors" {
		t.Errorf("Children[0] = %q, want %q", scenarios.Children[0], "using-k6/scenarios/executors")
	}

	// using-k6/scenarios/executors should have shared-iterations.
	executors := sections[4]
	if len(executors.Children) != 1 {
		t.Fatalf("executors children: got %d, want 1", len(executors.Children))
	}

	// Non-index leaf nodes should have empty (non-nil) children.
	checks := sections[1]
	if checks.Children == nil {
		t.Error("leaf Children should be non-nil")
	}
	if len(checks.Children) != 0 {
		t.Errorf("leaf Children count = %d, want 0", len(checks.Children))
	}
}

func TestSlugCollisionPrefersIndex(t *testing.T) {
	t.Parallel()

	version := "v0.99.x"
	root := t.TempDir()
	versionRoot := filepath.Join(root, "docs", "sources", "k6", version)

	// Create a regular file and an _index.md that produce the same slug.
	// javascript-api/k6-http/cookiejar.md  -> slug: javascript-api/k6-http/cookiejar
	// javascript-api/k6-http/cookiejar/_index.md -> slug: javascript-api/k6-http/cookiejar
	writeFile(t, filepath.Join(versionRoot, "_index.md"), "---\ntitle: root\n---\n")
	writeFile(t, filepath.Join(versionRoot, "javascript-api", "_index.md"), "---\ntitle: 'JS API'\nweight: 1\n---\n")
	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "_index.md"), "---\ntitle: 'k6/http'\nweight: 1\n---\n")
	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "cookiejar.md"),
		"---\ntitle: 'cookiejar function'\nweight: 10\n---\n\nA function.\n")
	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "cookiejar", "_index.md"),
		"---\ntitle: 'CookieJar class'\nweight: 20\n---\n\nA class with children.\n")
	writeFile(t, filepath.Join(versionRoot, "javascript-api", "k6-http", "cookiejar", "set.md"),
		"---\ntitle: 'set'\nweight: 1\n---\n\nSet a cookie.\n")

	outputDir := filepath.Join(t.TempDir(), "output")
	if err := run(version, root, outputDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "sections.json"))
	if err != nil {
		t.Fatalf("read sections.json: %v", err)
	}

	var idx docs.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parse sections.json: %v", err)
	}

	// Count how many sections have the colliding slug.
	slug := "javascript-api/k6-http/cookiejar"
	var matches []docs.Section
	for _, s := range idx.Sections {
		if s.Slug == slug {
			matches = append(matches, s)
		}
	}

	if len(matches) != 1 {
		t.Fatalf("expected exactly 1 section with slug %q, got %d", slug, len(matches))
	}

	// The _index.md version should win (it has children).
	if !matches[0].IsIndex {
		t.Errorf("expected the _index.md version to win the slug collision, got IsIndex=false")
	}
	if matches[0].Title != "CookieJar class" {
		t.Errorf("Title = %q, want %q", matches[0].Title, "CookieJar class")
	}
}

func TestRunWithRealDocs(t *testing.T) {
	k6DocsPath := os.Getenv("K6_DOCS_PATH")
	if k6DocsPath == "" {
		k6DocsPath = filepath.Join(os.Getenv("HOME"), "grafana", "k6-docs")
	}
	if _, err := os.Stat(k6DocsPath); err != nil {
		t.Skipf("skipping integration test: k6-docs not found at %s", k6DocsPath)
	}

	outputDir := filepath.Join(t.TempDir(), "real-output")
	version := "v1.5.x"

	if err := run(version, k6DocsPath, outputDir); err != nil {
		t.Fatalf("run with real docs: %v", err)
	}

	// Validate sections.json.
	data, err := os.ReadFile(filepath.Join(outputDir, "sections.json"))
	if err != nil {
		t.Fatalf("read sections.json: %v", err)
	}

	var idx docs.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parse sections.json: %v", err)
	}

	if idx.Version != version {
		t.Errorf("Version = %q, want %q", idx.Version, version)
	}

	// Should have a reasonable number of sections.
	if len(idx.Sections) < 50 {
		t.Errorf("expected at least 50 sections, got %d", len(idx.Sections))
	}

	// Build slug map.
	bySlug := make(map[string]docs.Section, len(idx.Sections))
	for _, s := range idx.Sections {
		bySlug[s.Slug] = s
	}

	// Spot-check some expected sections.
	spotChecks := []string{
		"javascript-api",
		"javascript-api/k6-http",
		"javascript-api/k6-http/get",
		"using-k6",
		"using-k6/checks",
		"reference/glossary",
	}
	for _, slug := range spotChecks {
		if _, ok := bySlug[slug]; !ok {
			t.Errorf("expected section %q in real output", slug)
		}
	}

	// Verify excluded categories are absent.
	for _, s := range idx.Sections {
		cat := s.Category
		excluded := map[string]bool{
			"get-started":      true,
			"set-up":           true,
			"extensions":       true,
			"grafana-cloud-k6": true,
			"release-notes":    true,
			"k6-studio":        true,
		}
		if excluded[cat] {
			t.Errorf("section %q has excluded category %q", s.Slug, cat)
		}
	}

	// Verify reference category only has glossary.
	for _, s := range idx.Sections {
		if s.Category == "reference" && s.Slug != "reference/glossary" {
			t.Errorf("reference should only include glossary, found %q", s.Slug)
		}
	}

	// Check that a transformed markdown file exists and has no frontmatter.
	checksPath := filepath.Join(outputDir, "markdown", "using-k6", "checks.md")
	checksData, err := os.ReadFile(checksPath)
	if err != nil {
		t.Fatalf("read transformed checks.md: %v", err)
	}
	if strings.HasPrefix(string(checksData), "---") {
		t.Error("transformed markdown should not start with frontmatter")
	}
	if strings.Contains(string(checksData), "{{<") {
		t.Error("transformed markdown should not contain Hugo shortcodes")
	}

	// best_practices.md should exist.
	if _, err := os.Stat(filepath.Join(outputDir, "best_practices.md")); err != nil {
		t.Error("best_practices.md should exist in output")
	}
}

func TestRunWithExactVersionNoVPrefix(t *testing.T) {
	t.Parallel()

	// The docs directory uses the wildcard form v0.99.x. When the caller
	// passes an exact version without the "v" prefix (e.g. "0.99.3"),
	// MapToWildcard must still produce "v0.99.x" to match the directory.
	docsPath := setupMockDocs(t, "v0.99.x")
	outputDir := filepath.Join(t.TempDir(), "output")

	if err := run("0.99.3", docsPath, outputDir); err != nil {
		t.Fatalf("run with bare version (no v prefix): %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "sections.json"))
	if err != nil {
		t.Fatalf("read sections.json: %v", err)
	}

	var idx docs.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parse sections.json: %v", err)
	}

	if idx.Version != "0.99.3" {
		t.Errorf("Version = %q, want %q (original version should be preserved)", idx.Version, "0.99.3")
	}

	if len(idx.Sections) == 0 {
		t.Error("expected sections to be populated")
	}
}

func TestMissingVersion(t *testing.T) {
	t.Parallel()

	docsPath := setupMockDocs(t, "v0.99.x")
	outputDir := filepath.Join(t.TempDir(), "output")

	err := run("v999.999.x", docsPath, outputDir)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
	if !strings.Contains(err.Error(), "version root not found") {
		t.Errorf("error = %q, expected to mention 'version root not found'", err.Error())
	}
}

func TestRunWithExactVersion(t *testing.T) {
	t.Parallel()

	// The docs directory uses the wildcard form v0.99.x, but the caller
	// passes an exact version like v0.99.3. The run function must map the
	// exact version to the wildcard directory automatically.
	docsPath := setupMockDocs(t, "v0.99.x")
	outputDir := filepath.Join(t.TempDir(), "output")

	if err := run("v0.99.3", docsPath, outputDir); err != nil {
		t.Fatalf("run with exact version: %v", err)
	}

	// The index should preserve the original exact version.
	data, err := os.ReadFile(filepath.Join(outputDir, "sections.json"))
	if err != nil {
		t.Fatalf("read sections.json: %v", err)
	}

	var idx docs.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parse sections.json: %v", err)
	}

	if idx.Version != "v0.99.3" {
		t.Errorf("Version = %q, want %q (exact version should be preserved)", idx.Version, "v0.99.3")
	}

	if len(idx.Sections) == 0 {
		t.Error("expected sections to be populated")
	}
}

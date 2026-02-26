package docs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// slugToShortArgs converts a canonical slug to the short CLI form.
//
// Examples:
//
//	"javascript-api/k6-http/get"           -> "http get"
//	"javascript-api/k6-browser/page/click" -> "browser page click"
//	"using-k6/scenarios"                   -> "using-k6 scenarios"
//	"examples/websockets"                  -> "examples websockets"
func slugToShortArgs(slug string) string {
	parts := strings.Split(slug, "/")
	if len(parts) == 0 {
		return slug
	}

	// JavaScript API modules: strip "javascript-api" prefix and "k6-" from the module name.
	if parts[0] == "javascript-api" && len(parts) > 1 {
		module := strings.TrimPrefix(parts[1], "k6-")
		rest := parts[2:]
		out := append([]string{module}, rest...)
		return strings.Join(out, " ")
	}

	// Everything else: join parts with spaces.
	return strings.Join(parts, " ")
}

// printTOC prints the table of contents grouped by category.
func printTOC(w io.Writer, idx *Index, version string) {
	fmt.Fprintf(w, "k6 Documentation (%s)\n\n", version)
	fmt.Fprintln(w, "Use: k6 x docs <topic>")

	topLevel := idx.TopLevel()

	for _, cat := range topLevel {
		fmt.Fprintf(w, "\n## %s\n", cat.Title)

		children := idx.Children(cat.Slug)
		if len(children) == 0 {
			// Show the category itself if it has no children.
			fmt.Fprintf(w, "  %-20s %s\n", slugToShortArgs(cat.Slug), cat.Description)
			continue
		}

		for _, child := range children {
			name := slugToShortArgs(child.Slug)
			fmt.Fprintf(w, "  %-20s %s\n", name, child.Description)
		}
	}
}

// printSection prints a section's markdown content, read from the cache dir.
// If the section has children, a subtopics footer is appended.
func printSection(w io.Writer, idx *Index, section *Section, cacheDir, version string) {
	content := readMarkdown(cacheDir, section.RelPath)
	if content != "" {
		content = Transform(content, version, nil)
		fmt.Fprint(w, content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Fprintln(w)
		}
	}

	children := idx.Children(section.Slug)
	if len(children) > 0 {
		names := make([]string, 0, len(children))
		for _, c := range children {
			names = append(names, slugToShortArgs(c.Slug))
		}

		fmt.Fprintln(w, "---")
		fmt.Fprintf(w, "Subtopics: %s\n", strings.Join(names, ", "))
		fmt.Fprintf(w, "Use: k6 x docs %s <subtopic>\n", slugToShortArgs(section.Slug))
	}
}

// printList prints children of a section in compact format.
func printList(w io.Writer, idx *Index, slug string) {
	sec, ok := idx.Lookup(slug)
	if !ok {
		fmt.Fprintf(w, "Topic not found: %s\n", slug)
		return
	}

	children := idx.Children(slug)

	fmt.Fprintf(w, "%s", sec.Title)
	if sec.Description != "" {
		fmt.Fprintf(w, " — %s", sec.Description)
	}
	fmt.Fprintln(w)

	if len(children) == 0 {
		fmt.Fprintln(w, "\n  (no subtopics)")
		return
	}

	fmt.Fprintln(w)
	for _, child := range children {
		name := slugToShortArgs(child.Slug)
		fmt.Fprintf(w, "  %-20s %s\n", name, child.Description)
	}
}

// printSearch prints search results formatted as usable CLI paths.
func printSearch(w io.Writer, idx *Index, term, cacheDir string) {
	readContent := func(slug string) string {
		sec, ok := idx.Lookup(slug)
		if !ok {
			return ""
		}
		return readMarkdown(cacheDir, sec.RelPath)
	}

	results := idx.Search(term, readContent)

	fmt.Fprintf(w, "Results for %q:\n", term)

	if len(results) == 0 {
		fmt.Fprintln(w, "\n  (no results)")
		return
	}

	fmt.Fprintln(w)
	for _, sec := range results {
		name := slugToShortArgs(sec.Slug)
		fmt.Fprintf(w, "  %-20s %s\n", name, sec.Description)
	}
}

// printBestPractices reads and prints the best_practices.md file from the cache.
func printBestPractices(w io.Writer, cacheDir string) error {
	path := filepath.Join(cacheDir, "best_practices.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read best practices: %w", err)
	}
	content := StripFrontmatter(string(data))
	fmt.Fprint(w, content)
	if !strings.HasSuffix(content, "\n") {
		fmt.Fprintln(w)
	}
	return nil
}

// printAll prints all sections sequentially.
func printAll(w io.Writer, idx *Index, cacheDir, version string) {
	fmt.Fprintf(w, "k6 Documentation (%s)\n\n", version)

	for i := range idx.Sections {
		sec := &idx.Sections[i]
		content := readMarkdown(cacheDir, sec.RelPath)
		if content == "" {
			continue
		}
		content = Transform(content, version, nil)

		fmt.Fprintf(w, "# %s\n\n", sec.Title)
		fmt.Fprint(w, content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w)
	}
}

// readMarkdown reads a markdown file from the cache directory.
func readMarkdown(cacheDir, relPath string) string {
	path := filepath.Join(cacheDir, "markdown", relPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

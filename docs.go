package docs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// childName returns the short name of a child relative to its parent.
// If the child slug starts with parentSlug+"/", the prefix is stripped.
// Otherwise, the last path segment is returned.
func childName(childSlug, parentSlug string) string {
	if strings.HasPrefix(childSlug, parentSlug+"/") {
		return childSlug[len(parentSlug)+1:]
	}
	if i := strings.LastIndex(childSlug, "/"); i >= 0 {
		return childSlug[i+1:]
	}
	return childSlug
}

// truncate shortens s to max characters, appending "..." if truncated.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// printTOC prints the table of contents grouped by category.
func printTOC(w io.Writer, idx *Index, version string) {
	fmt.Fprintf(w, "k6 Documentation (%s)\n", version)
	fmt.Fprintln(w, "Use: k6 x docs <topic>")

	topLevel := idx.TopLevel()

	for _, cat := range topLevel {
		fmt.Fprintf(w, "\n## %s\n", cat.Title)

		children := idx.Children(cat.Slug)
		if len(children) == 0 {
			// Show the category itself if it has no children.
			fmt.Fprintf(w, "  %-20s %s\n", childName(cat.Slug, ""), truncate(cat.Description, 80))
			continue
		}

		for _, child := range children {
			name := childName(child.Slug, cat.Slug)
			fmt.Fprintf(w, "  %-20s %s\n", name, truncate(child.Description, 80))
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
			names = append(names, childName(c.Slug, section.Slug))
		}

		fmt.Fprintln(w)
		fmt.Fprintln(w, "---")
		fmt.Fprintf(w, "Subtopics: %s\n", strings.Join(names, ", "))
		fmt.Fprintf(w, "Use: k6 x docs %s <subtopic>\n", childName(section.Slug, ""))
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

	for _, child := range children {
		name := childName(child.Slug, slug)
		fmt.Fprintf(w, "  %-20s %s\n", name, truncate(child.Description, 80))
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

	for _, sec := range results {
		name := childName(sec.Slug, "")
		fmt.Fprintf(w, "  %-20s %s\n", name, truncate(sec.Description, 80))
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

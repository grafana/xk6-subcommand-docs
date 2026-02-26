package docs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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

// searchGroupKey returns the grouping key for a search result.
// JavaScript API sections group by module (second segment); others by first segment.
func searchGroupKey(slug string) string {
	parts := strings.SplitN(slug, "/", 3)
	if parts[0] == "javascript-api" && len(parts) > 1 {
		return parts[1]
	}
	return parts[0]
}

// printSearch prints search results grouped hierarchically by topic.
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

	// Build a set of matched slugs and group results by parent topic.
	matched := make(map[string]*Section, len(results))
	groups := make(map[string][]*Section)
	var groupOrder []string

	for _, sec := range results {
		matched[sec.Slug] = sec

		// Skip bare "javascript-api" — it's the TOC default.
		if sec.Slug == "javascript-api" {
			continue
		}

		key := searchGroupKey(sec.Slug)
		if _, exists := groups[key]; !exists {
			groupOrder = append(groupOrder, key)
		}
		groups[key] = append(groups[key], sec)
	}

	// Sort groups alphabetically.
	sort.Strings(groupOrder)

	for _, key := range groupOrder {
		members := groups[key]

		// Sort items within group alphabetically by slug.
		sort.Slice(members, func(i, j int) bool {
			return members[i].Slug < members[j].Slug
		})

		// Check if the group topic itself is a matched result.
		// For JS API modules, the group slug is "javascript-api/{key}".
		// For others, it's just "{key}".
		groupSlug := key
		if _, ok := idx.Lookup("javascript-api/" + key); ok {
			// If there's a javascript-api/{key} section, this is a JS API module group.
			if members[0].Slug == "javascript-api/"+key || strings.HasPrefix(members[0].Slug, "javascript-api/"+key+"/") {
				groupSlug = "javascript-api/" + key
			}
		}

		groupSec := matched[groupSlug]

		// Print group header.
		if groupSec != nil {
			fmt.Fprintf(w, "%s: %s\n", key, truncate(groupSec.Description, 80))
		} else {
			fmt.Fprintf(w, "%s:\n", key)
		}

		// Print children (items that aren't the group header itself).
		for _, sec := range members {
			if sec.Slug == groupSlug {
				continue
			}
			name := childName(sec.Slug, groupSlug)
			fmt.Fprintf(w, "  %-22s %s\n", name, truncate(sec.Description, 80))
		}

		fmt.Fprintln(w)
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

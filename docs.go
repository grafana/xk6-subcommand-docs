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

// listItem is a name+description pair for aligned list rendering.
type listItem struct {
	Name        string
	Description string
}

// printAlignedList prints items as a left-aligned name+description list.
// Duplicate names (by Name) are skipped. Each line is indented by indent.
func printAlignedList(w io.Writer, items []listItem, indent string) {
	seen := make(map[string]bool, len(items))

	maxWidth := 0
	for _, item := range items {
		if seen[item.Name] {
			continue
		}
		seen[item.Name] = true
		if len(item.Name) > maxWidth {
			maxWidth = len(item.Name)
		}
	}

	fmtStr := fmt.Sprintf("%s%%-%ds %%s\n", indent, maxWidth+1)

	// Reset seen for the printing pass.
	for k := range seen {
		delete(seen, k)
	}
	for _, item := range items {
		if seen[item.Name] {
			continue
		}
		seen[item.Name] = true
		fmt.Fprintf(w, fmtStr, item.Name, truncate(item.Description, 80))
	}
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
			fmt.Fprintf(w, "- %s %s\n", childName(cat.Slug, ""), truncate(cat.Description, 80))
			continue
		}

		items := make([]listItem, 0, len(children))
		for _, child := range children {
			items = append(items, listItem{
				Name:        childName(child.Slug, cat.Slug),
				Description: child.Description,
			})
		}
		printAlignedList(w, items, "- ")
		fmt.Fprintf(w, "\n  → Usage: k6 x docs %s <topic>\n", cat.Slug)
	}
}

// printSection prints a section's markdown content, read from the cache dir.
// If the section has children, a subtopics footer is appended.
func printSection(w io.Writer, idx *Index, section *Section, cacheDir, version string) {
	content := readAndTransform(cacheDir, section.RelPath, version)
	if content != "" {
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

	items := make([]listItem, 0, len(children))
	for _, child := range children {
		items = append(items, listItem{
			Name:        childName(child.Slug, slug),
			Description: child.Description,
		})
	}
	printAlignedList(w, items, "- ")
}

// printTopLevelList lists all top-level categories with their descriptions.
func printTopLevelList(w io.Writer, idx *Index) {
	cats := idx.TopLevel()

	items := make([]listItem, 0, len(cats))
	for _, cat := range cats {
		items = append(items, listItem{
			Name:        cat.Slug,
			Description: cat.Description,
		})
	}
	printAlignedList(w, items, "- ")
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
func printSearch(w io.Writer, idx *Index, term, cacheDir, version string) {
	readContent := func(slug string) string {
		sec, ok := idx.Lookup(slug)
		if !ok {
			return ""
		}
		return readAndTransform(cacheDir, sec.RelPath, version)
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

		// Collect children (items that aren't the group header itself).
		var items []listItem
		for _, sec := range members {
			if sec.Slug == groupSlug {
				continue
			}
			items = append(items, listItem{
				Name:        childName(sec.Slug, groupSlug),
				Description: sec.Description,
			})
		}
		printAlignedList(w, items, "- ")

		fmt.Fprintln(w)
	}
}

// printBestPractices reads and prints the best_practices.md file from the cache.
func printBestPractices(w io.Writer, cacheDir, version string) error {
	path := filepath.Join(cacheDir, "best_practices.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read best practices: %w", err)
	}
	content := Transform(string(data), version)
	fmt.Fprint(w, content)
	if !strings.HasSuffix(content, "\n") {
		fmt.Fprintln(w)
	}
	return nil
}

// printAll prints all sections sequentially.
func printAll(w io.Writer, idx *Index, cacheDir, version string) {
	fmt.Fprintf(w, "k6 Documentation (%s)\n", version)

	for i := range idx.Sections {
		sec := &idx.Sections[i]
		content := readAndTransform(cacheDir, sec.RelPath, version)
		if content == "" {
			continue
		}
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

// readAndTransform reads a markdown file and applies runtime transforms.
func readAndTransform(cacheDir, relPath, version string) string {
	raw := readMarkdown(cacheDir, relPath)
	if raw == "" {
		return ""
	}
	return Transform(raw, version)
}

// Command prepare processes the k6-docs repository into a doc bundle
// suitable for embedding. It walks the documentation tree, transforms
// Hugo shortcodes into clean markdown, and produces:
//   - markdown/ — transformed .md files
//   - sections.json — structured index of all sections
//   - best_practices.md — a comprehensive best practices guide
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	docs "github.com/grafana/xk6-subcommand-docs"
	"gopkg.in/yaml.v3"
)

// includedCategories lists the top-level directories we keep.
var includedCategories = map[string]bool{
	"javascript-api":   true,
	"using-k6":         true,
	"using-k6-browser": true,
	"testing-guides":   true,
	"examples":         true,
	"results-output":   true,
}

// frontmatter holds the YAML fields we extract from each doc file.
type frontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Weight      int    `yaml:"weight"`
}

func main() {
	log.SetFlags(0)

	var (
		k6Version  string
		k6DocsPath string
		outputDir  string
	)

	flag.StringVar(&k6Version, "k6-version", "", "k6 docs version (e.g. v1.5.x) — required")
	flag.StringVar(&k6DocsPath, "k6-docs-path", "", "local path to k6-docs repo (cloned if empty)")
	flag.StringVar(&outputDir, "output-dir", "dist/", "output directory")
	flag.Parse()

	if k6Version == "" {
		log.Fatal("--k6-version is required")
	}

	if err := run(k6Version, k6DocsPath, outputDir); err != nil {
		log.Fatal(err)
	}
}

func run(k6Version, k6DocsPath, outputDir string) error {
	// Step 1: ensure we have the k6-docs repo.
	docsPath, cleanup, err := ensureDocsRepo(k6DocsPath)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	// The k6-docs repo uses wildcard directories (e.g. "v1.6.x"), so convert
	// exact versions like "v1.6.1" to the wildcard form for the path lookup.
	docsVersion := docs.MapToWildcard(k6Version)
	versionRoot := filepath.Join(docsPath, "docs", "sources", "k6", docsVersion)
	if _, err := os.Stat(versionRoot); err != nil {
		return fmt.Errorf("version root not found: %w", err)
	}

	// Step 2: build shared content map.
	sharedContent, err := buildSharedContentMap(filepath.Join(versionRoot, "shared"))
	if err != nil {
		return fmt.Errorf("build shared content: %w", err)
	}

	// Step 3: walk documentation files and collect sections.
	markdownDir := filepath.Join(outputDir, "markdown")
	sections, err := walkAndProcess(versionRoot, markdownDir, k6Version, sharedContent)
	if err != nil {
		return fmt.Errorf("walk docs: %w", err)
	}

	// Step 4: populate children.
	populateChildren(sections)

	// Step 5: write sections.json.
	idx := docs.Index{
		Version:  k6Version,
		Sections: sections,
	}
	if err := writeSectionsJSON(outputDir, idx); err != nil {
		return err
	}

	// Step 6: write best_practices.md.
	if err := writeBestPractices(outputDir); err != nil {
		return err
	}

	log.Printf("Done: %d sections written to %s", len(sections), outputDir)
	return nil
}

// ensureDocsRepo returns the path to the k6-docs repo. If k6DocsPath is empty,
// it clones the repo into a temp directory and returns a cleanup function.
func ensureDocsRepo(k6DocsPath string) (string, func(), error) {
	if k6DocsPath != "" {
		return k6DocsPath, nil, nil
	}

	tmpDir, err := os.MkdirTemp("", "k6-docs-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp dir: %w", err)
	}

	log.Println("Cloning k6-docs repository...")
	cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/grafana/k6-docs.git", tmpDir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("clone k6-docs: %w", err)
	}

	cleanup := func() { os.RemoveAll(tmpDir) }
	return tmpDir, cleanup, nil
}

// buildSharedContentMap reads all .md files under the shared directory and
// returns a map keyed by the relative path (e.g. "javascript-api/k6-http.md").
func buildSharedContentMap(sharedDir string) (map[string]string, error) {
	m := make(map[string]string)

	info, err := os.Stat(sharedDir)
	if err != nil || !info.IsDir() {
		return m, nil
	}

	err = filepath.WalkDir(sharedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, err := filepath.Rel(sharedDir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read shared %s: %w", rel, err)
		}
		m[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	return m, err
}

// isIncluded reports whether a relative path from the version root should be included.
func isIncluded(relPath string) bool {
	// Normalize to forward slashes for consistent matching.
	relPath = filepath.ToSlash(relPath)

	parts := strings.SplitN(relPath, "/", 2)
	if len(parts) == 0 {
		return false
	}
	topDir := parts[0]

	// Direct category match.
	if includedCategories[topDir] {
		return true
	}

	// Special case: reference/glossary only.
	if topDir == "reference" {
		return strings.HasPrefix(relPath, "reference/glossary")
	}

	return false
}

// parseFrontmatter extracts YAML frontmatter from content.
func parseFrontmatter(content string) (frontmatter, error) {
	var fm frontmatter
	if !strings.HasPrefix(content, "---\n") {
		return fm, nil
	}
	end := strings.Index(content[4:], "\n---")
	if end == -1 {
		return fm, nil
	}
	yamlBlock := deduplicateYAMLKeys(content[4 : 4+end])
	if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err != nil {
		return fm, fmt.Errorf("parse yaml: %w", err)
	}
	return fm, nil
}

// deduplicateYAMLKeys removes duplicate top-level YAML keys, keeping only
// the first occurrence of each key. This handles the ~60 k6-docs files that
// have duplicate "description:" keys, which cause yaml.v3 to error.
func deduplicateYAMLKeys(yamlBlock string) string {
	seen := make(map[string]bool)
	var lines []string
	for _, line := range strings.Split(yamlBlock, "\n") {
		if idx := strings.Index(line, ":"); idx > 0 && len(line) > 0 && line[0] != ' ' && line[0] != '\t' && line[0] != '#' {
			key := strings.TrimSpace(line[:idx])
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// slugFromRelPath derives the slug from a relative path.
// Rules: strip .md, if _index.md use parent dir, path uses forward slashes.
func slugFromRelPath(relPath string) string {
	relPath = filepath.ToSlash(relPath)
	base := filepath.Base(relPath)
	if base == "_index.md" {
		return filepath.ToSlash(filepath.Dir(relPath))
	}
	return strings.TrimSuffix(relPath, ".md")
}

// categoryFromSlug extracts the first path segment as the category.
func categoryFromSlug(slug string) string {
	if i := strings.Index(slug, "/"); i != -1 {
		return slug[:i]
	}
	return slug
}

// walkAndProcess walks the version root, processes included .md files,
// and returns the collected sections.
func walkAndProcess(versionRoot, markdownDir, k6Version string, sharedContent map[string]string) ([]docs.Section, error) {
	// Use a map to deduplicate sections by slug. When a slug collision
	// occurs (e.g. cookiejar.md and cookiejar/_index.md both produce
	// "javascript-api/k6-http/cookiejar"), prefer the _index.md entry
	// because it represents a section with children.
	sectionMap := make(map[string]docs.Section)
	var slugOrder []string

	err := filepath.WalkDir(versionRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(versionRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		// Skip the shared directory entirely.
		if d.IsDir() && rel == "shared" {
			return filepath.SkipDir
		}

		// Skip non-markdown files and directories.
		if d.IsDir() || !strings.HasSuffix(rel, ".md") {
			return nil
		}

		// Skip the version root _index.md.
		if rel == "_index.md" {
			return nil
		}

		// Only include files from allowed categories.
		if !isIncluded(rel) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}

		fm, err := parseFrontmatter(string(content))
		if err != nil {
			log.Printf("warning: %s: %v", rel, err)
		}

		transformed := docs.Transform(string(content), k6Version, sharedContent)

		slug := slugFromRelPath(rel)
		category := categoryFromSlug(slug)
		isIndex := filepath.Base(path) == "_index.md"

		// Write transformed markdown.
		outPath := filepath.Join(markdownDir, rel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
		}
		if err := os.WriteFile(outPath, []byte(transformed), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}

		sec := docs.Section{
			Slug:        slug,
			RelPath:     rel,
			Title:       fm.Title,
			Description: fm.Description,
			Weight:      fm.Weight,
			Category:    category,
			IsIndex:     isIndex,
		}

		// Handle slug collisions: prefer _index.md over plain .md files.
		if existing, ok := sectionMap[slug]; ok {
			if isIndex && !existing.IsIndex {
				sectionMap[slug] = sec
			}
			// Otherwise keep the existing entry.
		} else {
			slugOrder = append(slugOrder, slug)
			sectionMap[slug] = sec
		}

		return nil
	})

	// Rebuild the slice in walk order.
	sections := make([]docs.Section, 0, len(slugOrder))
	for _, slug := range slugOrder {
		sections = append(sections, sectionMap[slug])
	}

	return sections, err
}

// populateChildren sets the Children field for each _index section.
// A child is a section whose slug starts with parent slug + "/" and has
// no further "/" after that prefix (direct child only).
func populateChildren(sections []docs.Section) {
	for i := range sections {
		if !sections[i].IsIndex {
			continue
		}

		parentSlug := sections[i].Slug
		prefix := parentSlug + "/"

		// Collect direct children.
		type child struct {
			slug   string
			weight int
		}
		var children []child

		for j := range sections {
			if i == j {
				continue
			}
			s := sections[j].Slug
			if !strings.HasPrefix(s, prefix) {
				continue
			}
			remainder := s[len(prefix):]
			if strings.Contains(remainder, "/") {
				continue
			}
			children = append(children, child{slug: s, weight: sections[j].Weight})
		}

		sort.Slice(children, func(a, b int) bool {
			return children[a].weight < children[b].weight
		})

		slugs := make([]string, len(children))
		for k, c := range children {
			slugs[k] = c.slug
		}
		sections[i].Children = slugs
	}

	// Ensure non-index sections have empty (non-nil) Children.
	for i := range sections {
		if sections[i].Children == nil {
			sections[i].Children = []string{}
		}
	}
}

// writeSectionsJSON writes the index to sections.json in the output directory.
func writeSectionsJSON(outputDir string, idx docs.Index) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal sections: %w", err)
	}

	outPath := filepath.Join(outputDir, "sections.json")
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return fmt.Errorf("write sections.json: %w", err)
	}

	log.Printf("Wrote %s", outPath)
	return nil
}

// writeBestPractices writes a comprehensive best practices guide.
func writeBestPractices(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	outPath := filepath.Join(outputDir, "best_practices.md")
	if err := os.WriteFile(outPath, []byte(bestPracticesContent), 0o644); err != nil {
		return fmt.Errorf("write best_practices.md: %w", err)
	}

	log.Printf("Wrote %s", outPath)
	return nil
}

//go:embed best_practices.md
var bestPracticesContent string

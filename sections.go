package docs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Section represents a single documentation section.
type Section struct {
	Slug        string   `json:"slug"`
	RelPath     string   `json:"rel_path"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Weight      int      `json:"weight"`
	Category    string   `json:"category"`
	Children    []string `json:"children"`
	IsIndex     bool     `json:"is_index"`
}

// Index holds all sections and provides fast lookup by slug.
type Index struct {
	Version  string    `json:"version"`
	Sections []Section `json:"sections"`
	bySlug   map[string]*Section
}

// LoadIndex reads sections.json from dir and returns a populated Index.
func LoadIndex(dir string) (*Index, error) {
	data, err := os.ReadFile(filepath.Join(dir, "sections.json"))
	if err != nil {
		return nil, fmt.Errorf("load index %s: %w", dir, err)
	}

	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse index %s: %w", dir, err)
	}

	idx.bySlug = make(map[string]*Section, len(idx.Sections))
	for i := range idx.Sections {
		idx.bySlug[idx.Sections[i].Slug] = &idx.Sections[i]
	}

	return &idx, nil
}

// Lookup returns the section with the given slug in O(1) time.
func (idx *Index) Lookup(slug string) (*Section, bool) {
	sec, ok := idx.bySlug[slug]
	return sec, ok
}

// Search returns sections whose title, description, or body (via readContent)
// contain term as a case-insensitive substring. If readContent is nil, only
// title and description are checked.
func (idx *Index) Search(term string, readContent func(slug string) string) []*Section {
	if term == "" {
		return nil
	}

	lower := strings.ToLower(term)
	var results []*Section

	for i := range idx.Sections {
		sec := &idx.Sections[i]

		if strings.Contains(strings.ToLower(sec.Title), lower) ||
			strings.Contains(strings.ToLower(sec.Description), lower) {
			results = append(results, sec)
			continue
		}

		if readContent != nil {
			body := readContent(sec.Slug)
			if body != "" && strings.Contains(strings.ToLower(body), lower) {
				results = append(results, sec)
			}
		}
	}

	return results
}

// Children returns the child sections of the given slug, sorted by weight.
// Returns nil if the slug is not found.
func (idx *Index) Children(slug string) []*Section {
	parent, ok := idx.bySlug[slug]
	if !ok {
		return nil
	}

	children := make([]*Section, 0, len(parent.Children))
	for _, childSlug := range parent.Children {
		if child, ok := idx.bySlug[childSlug]; ok {
			children = append(children, child)
		}
	}

	sort.Slice(children, func(i, j int) bool {
		return children[i].Weight < children[j].Weight
	})

	return children
}

// TopLevel returns sections where Category == Slug (top-level indices),
// sorted by weight.
func (idx *Index) TopLevel() []*Section {
	var top []*Section
	for i := range idx.Sections {
		sec := &idx.Sections[i]
		if sec.Category == sec.Slug {
			top = append(top, sec)
		}
	}

	sort.Slice(top, func(i, j int) bool {
		return top[i].Weight < top[j].Weight
	})

	return top
}

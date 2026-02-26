package docs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadIndex(t *testing.T) {
	t.Run("valid fixture", func(t *testing.T) {
		idx, err := LoadIndex("testdata")
		if err != nil {
			t.Fatalf("LoadIndex: unexpected error: %v", err)
		}
		if idx.Version != "1.0.0" {
			t.Errorf("Version = %q, want %q", idx.Version, "1.0.0")
		}
		if got := len(idx.Sections); got != 9 {
			t.Errorf("len(Sections) = %d, want 9", got)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := LoadIndex("/tmp/nonexistent-dir-xk6-test")
		if err == nil {
			t.Fatal("LoadIndex: expected error for missing directory, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "sections.json"), []byte("{bad json"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := LoadIndex(dir)
		if err == nil {
			t.Fatal("LoadIndex: expected error for invalid JSON, got nil")
		}
	})
}

func TestLookup(t *testing.T) {
	idx := mustLoadIndex(t)

	t.Run("existing slug", func(t *testing.T) {
		sec, ok := idx.Lookup("installation")
		if !ok {
			t.Fatal("Lookup(installation): not found")
		}
		if sec.Title != "Installation" {
			t.Errorf("Title = %q, want %q", sec.Title, "Installation")
		}
		if sec.Category != "getting-started" {
			t.Errorf("Category = %q, want %q", sec.Category, "getting-started")
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		sec, ok := idx.Lookup("Installation")
		if !ok {
			t.Fatal("Lookup(Installation): not found")
		}
		if sec.Slug != "installation" {
			t.Errorf("Slug = %q, want %q", sec.Slug, "installation")
		}
	})

	t.Run("missing slug", func(t *testing.T) {
		_, ok := idx.Lookup("does-not-exist")
		if ok {
			t.Error("Lookup(does-not-exist): expected false, got true")
		}
	})

	t.Run("empty slug", func(t *testing.T) {
		_, ok := idx.Lookup("")
		if ok {
			t.Error("Lookup(''): expected false, got true")
		}
	})
}

func TestSearch(t *testing.T) {
	idx := mustLoadIndex(t)

	// readContent simulates body content for specific slugs.
	readContent := func(slug string) string {
		if slug == "grpc" {
			return "gRPC uses protocol buffers for serialization."
		}
		return ""
	}

	t.Run("match in title", func(t *testing.T) {
		results := idx.Search("installation", nil)
		if len(results) != 1 {
			t.Fatalf("Search(installation): got %d results, want 1", len(results))
		}
		if results[0].Slug != "installation" {
			t.Errorf("Slug = %q, want %q", results[0].Slug, "installation")
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		results := idx.Search("GETTING STARTED", nil)
		if len(results) == 0 {
			t.Fatal("Search(GETTING STARTED): expected matches, got none")
		}
		foundParent := false
		for _, s := range results {
			if s.Slug == "getting-started" {
				foundParent = true
			}
		}
		if !foundParent {
			t.Error("Search: expected getting-started in results")
		}
	})

	t.Run("match in description", func(t *testing.T) {
		results := idx.Search("export test results", nil)
		if len(results) != 1 {
			t.Fatalf("Search(export test results): got %d results, want 1", len(results))
		}
		if results[0].Slug != "results" {
			t.Errorf("Slug = %q, want %q", results[0].Slug, "results")
		}
	})

	t.Run("match via readContent callback", func(t *testing.T) {
		results := idx.Search("protocol buffers", readContent)
		if len(results) != 1 {
			t.Fatalf("Search(protocol buffers): got %d results, want 1", len(results))
		}
		if results[0].Slug != "grpc" {
			t.Errorf("Slug = %q, want %q", results[0].Slug, "grpc")
		}
	})

	t.Run("no match", func(t *testing.T) {
		results := idx.Search("zzzznotfound", readContent)
		if len(results) != 0 {
			t.Errorf("Search(zzzznotfound): got %d results, want 0", len(results))
		}
	})

	t.Run("empty term", func(t *testing.T) {
		results := idx.Search("", nil)
		if len(results) != 0 {
			t.Errorf("Search(''): got %d results, want 0", len(results))
		}
	})

	t.Run("fuzzy: spaces match concatenated title", func(t *testing.T) {
		results := idx.Search("close context", nil)
		if len(results) != 1 {
			t.Fatalf("Search(close context): got %d results, want 1", len(results))
		}
		if results[0].Slug != "browser/closecontext" {
			t.Errorf("Slug = %q, want %q", results[0].Slug, "browser/closecontext")
		}
	})

	t.Run("fuzzy: dashes match concatenated title", func(t *testing.T) {
		results := idx.Search("close-context", nil)
		if len(results) != 1 {
			t.Fatalf("Search(close-context): got %d results, want 1", len(results))
		}
		if results[0].Slug != "browser/closecontext" {
			t.Errorf("Slug = %q, want %q", results[0].Slug, "browser/closecontext")
		}
	})

	t.Run("fuzzy: spaces match dashed slug", func(t *testing.T) {
		results := idx.Search("http debugging", nil)
		if len(results) != 1 {
			t.Fatalf("Search(http debugging): got %d results, want 1", len(results))
		}
		if results[0].Slug != "http-debugging" {
			t.Errorf("Slug = %q, want %q", results[0].Slug, "http-debugging")
		}
	})
}

func TestChildren(t *testing.T) {
	idx := mustLoadIndex(t)

	t.Run("parent with children", func(t *testing.T) {
		children := idx.Children("getting-started")
		if len(children) != 2 {
			t.Fatalf("Children(getting-started): got %d, want 2", len(children))
		}
		// Should be sorted by weight: installation (1), first-test (2)
		if children[0].Slug != "installation" {
			t.Errorf("children[0].Slug = %q, want %q", children[0].Slug, "installation")
		}
		if children[1].Slug != "first-test" {
			t.Errorf("children[1].Slug = %q, want %q", children[1].Slug, "first-test")
		}
	})

	t.Run("parent with no children", func(t *testing.T) {
		children := idx.Children("results")
		if len(children) != 0 {
			t.Errorf("Children(results): got %d, want 0", len(children))
		}
	})

	t.Run("nonexistent slug", func(t *testing.T) {
		children := idx.Children("nope")
		if children != nil {
			t.Errorf("Children(nope): expected nil, got %v", children)
		}
	})

	t.Run("protocols children sorted by weight", func(t *testing.T) {
		children := idx.Children("protocols")
		if len(children) != 2 {
			t.Fatalf("Children(protocols): got %d, want 2", len(children))
		}
		if children[0].Slug != "http" {
			t.Errorf("children[0].Slug = %q, want %q", children[0].Slug, "http")
		}
		if children[1].Slug != "grpc" {
			t.Errorf("children[1].Slug = %q, want %q", children[1].Slug, "grpc")
		}
	})
}

func TestTopLevel(t *testing.T) {
	idx := mustLoadIndex(t)

	top := idx.TopLevel()
	// Top-level: getting-started (w1), results (w2), protocols (w3)
	if len(top) != 3 {
		t.Fatalf("TopLevel: got %d, want 3", len(top))
	}

	// Verify sorted by weight
	if top[0].Slug != "getting-started" {
		t.Errorf("top[0].Slug = %q, want %q", top[0].Slug, "getting-started")
	}
	if top[1].Slug != "results" {
		t.Errorf("top[1].Slug = %q, want %q", top[1].Slug, "results")
	}
	if top[2].Slug != "protocols" {
		t.Errorf("top[2].Slug = %q, want %q", top[2].Slug, "protocols")
	}

	// All should be index sections
	for _, s := range top {
		if !s.IsIndex {
			t.Errorf("TopLevel section %q: IsIndex = false, want true", s.Slug)
		}
	}
}

func TestChildrenWithMissingChildSlug(t *testing.T) {
	// Test that Children gracefully skips slugs that don't exist in the index.
	idx := &Index{
		Sections: []Section{
			{Slug: "parent", Children: []string{"exists", "ghost"}, Weight: 1},
			{Slug: "exists", Weight: 1},
		},
	}
	idx.bySlug = map[string]*Section{
		"parent": &idx.Sections[0],
		"exists": &idx.Sections[1],
	}

	children := idx.Children("parent")
	if len(children) != 1 {
		t.Fatalf("Children(parent): got %d, want 1 (ghost should be skipped)", len(children))
	}
	if children[0].Slug != "exists" {
		t.Errorf("children[0].Slug = %q, want %q", children[0].Slug, "exists")
	}
}

// mustLoadIndex is a test helper that loads the fixture or fails the test.
func mustLoadIndex(t *testing.T) *Index {
	t.Helper()
	idx, err := LoadIndex("testdata")
	if err != nil {
		t.Fatalf("mustLoadIndex: %v", err)
	}
	return idx
}

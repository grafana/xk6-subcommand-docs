package docs

import (
	"archive/tar"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauspost/compress/zstd"
)

func TestDownloadURL(t *testing.T) {
	t.Parallel()

	got := downloadURL("v1.0.0")
	want := "https://github.com/grafana/xk6-subcommand-docs/releases/download/doc-bundles/docs-v1.0.0.tar.zst"
	if got != want {
		t.Errorf("downloadURL(v1.0.0) = %q, want %q", got, want)
	}
}

func TestCacheDir(t *testing.T) {
	t.Parallel()

	dir, err := CacheDir("v1.2.3")
	if err != nil {
		t.Fatalf("CacheDir: %v", err)
	}

	if !strings.HasSuffix(dir, filepath.Join("k6", "docs", "v1.2.3")) {
		t.Errorf("CacheDir = %q, want suffix %q", dir, filepath.Join("k6", "docs", "v1.2.3"))
	}
}

func TestIsCached(t *testing.T) {
	t.Parallel()

	// A version that definitely doesn't exist should not be cached.
	if IsCached("nonexistent-version-xyz") {
		t.Error("IsCached returned true for a version that should not exist")
	}

	// Create the directory and check again.
	dir, err := CacheDir("test-cached-version")
	if err != nil {
		t.Fatalf("CacheDir: %v", err)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	if !IsCached("test-cached-version") {
		t.Error("IsCached returned false after creating cache directory")
	}
}

func TestExtract(t *testing.T) {
	t.Parallel()

	// Build a .tar.zst archive in memory with two files.
	archive := buildTarZst(t, map[string]string{
		"readme.txt":        "hello world",
		"subdir/nested.txt": "nested content",
	})

	dest := t.TempDir()

	if err := extract(archive, dest); err != nil {
		t.Fatalf("extract: %v", err)
	}

	// Verify extracted files.
	assertFileContent(t, filepath.Join(dest, "readme.txt"), "hello world")
	assertFileContent(t, filepath.Join(dest, "subdir", "nested.txt"), "nested content")
}

func TestExtractRejectsTraversal(t *testing.T) {
	t.Parallel()

	archive := buildTarZstRaw(t, []tarEntry{
		{name: "../escape.txt", content: "evil"},
	})

	dest := t.TempDir()

	err := extract(archive, dest)
	if err == nil {
		t.Fatal("extract should reject path traversal, but returned nil")
	}
	if !strings.Contains(err.Error(), "traversal") {
		t.Errorf("expected traversal error, got: %v", err)
	}
}

func TestExtractRejectsAbsolutePath(t *testing.T) {
	t.Parallel()

	archive := buildTarZstRaw(t, []tarEntry{
		{name: "/etc/passwd", content: "evil"},
	})

	dest := t.TempDir()

	err := extract(archive, dest)
	if err == nil {
		t.Fatal("extract should reject absolute path, but returned nil")
	}
}

func TestEnsureDocs(t *testing.T) {
	t.Parallel()

	version := "test-ensure-" + t.Name()

	// Clean up any leftover cache from previous runs.
	dir, err := CacheDir(version)
	if err != nil {
		t.Fatalf("CacheDir: %v", err)
	}
	os.RemoveAll(dir)
	t.Cleanup(func() { os.RemoveAll(dir) })

	archive := buildTarZst(t, map[string]string{
		"doc.txt": "documentation content",
	})

	mock := &mockHTTPClient{
		body:       archive.Bytes(),
		statusCode: http.StatusOK,
	}

	got, err := EnsureDocs(version, mock)
	if err != nil {
		t.Fatalf("EnsureDocs: %v", err)
	}

	if got != dir {
		t.Errorf("EnsureDocs returned %q, want %q", got, dir)
	}

	assertFileContent(t, filepath.Join(dir, "doc.txt"), "documentation content")

	// Calling again should use cache (no second HTTP call).
	got2, err := EnsureDocs(version, mock)
	if err != nil {
		t.Fatalf("EnsureDocs second call: %v", err)
	}
	if got2 != dir {
		t.Errorf("second EnsureDocs returned %q, want %q", got2, dir)
	}
	if mock.calls != 1 {
		t.Errorf("expected 1 HTTP call, got %d", mock.calls)
	}
}

func TestEnsureDocsHTTPError(t *testing.T) {
	t.Parallel()

	version := "test-ensure-httperr-" + t.Name()
	dir, _ := CacheDir(version)
	os.RemoveAll(dir)
	t.Cleanup(func() { os.RemoveAll(dir) })

	mock := &mockHTTPClient{
		body:       []byte("not found"),
		statusCode: http.StatusNotFound,
	}

	_, err := EnsureDocs(version, mock)
	if err == nil {
		t.Fatal("EnsureDocs should fail on HTTP 404")
	}
}

// --- helpers ---

type mockHTTPClient struct {
	body       []byte
	statusCode int
	calls      int
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	m.calls++
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewReader(m.body)),
	}, nil
}

type tarEntry struct {
	name    string
	content string
}

func buildTarZst(t *testing.T, files map[string]string) *bytes.Buffer {
	t.Helper()

	entries := make([]tarEntry, 0, len(files))
	for name, content := range files {
		entries = append(entries, tarEntry{name: name, content: content})
	}
	return buildTarZstRaw(t, entries)
}

func buildTarZstRaw(t *testing.T, entries []tarEntry) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer

	zw, err := zstd.NewWriter(&buf)
	if err != nil {
		t.Fatalf("zstd.NewWriter: %v", err)
	}

	tw := tar.NewWriter(zw)
	for _, e := range entries {
		hdr := &tar.Header{
			Name: e.name,
			Mode: 0o644,
			Size: int64(len(e.content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("WriteHeader(%s): %v", e.name, err)
		}
		if _, err := tw.Write([]byte(e.content)); err != nil {
			t.Fatalf("Write(%s): %v", e.name, err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("tar.Close: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zstd.Close: %v", err)
	}

	return &buf
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	if got := string(data); got != want {
		t.Errorf("file %s content = %q, want %q", path, got, want)
	}
}

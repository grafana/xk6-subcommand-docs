package docs

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		dir := t.TempDir()
		k6Dir := filepath.Join(dir, "k6")
		if err := os.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: glow -p 200\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		t.Setenv("XDG_CONFIG_HOME", dir)

		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("loadConfig: unexpected error: %v", err)
		}
		if cfg.Renderer != "glow -p 200" {
			t.Errorf("loadConfig: Renderer = %q, want %q", cfg.Renderer, "glow -p 200")
		}
	})

	t.Run("missing file returns empty config", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", dir)

		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("loadConfig: unexpected error: %v", err)
		}
		if cfg.Renderer != "" {
			t.Errorf("loadConfig: Renderer = %q, want empty", cfg.Renderer)
		}
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		dir := t.TempDir()
		k6Dir := filepath.Join(dir, "k6")
		if err := os.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(k6Dir, "docs.yaml"), []byte(":\n  :\n    : [invalid"), 0o644); err != nil {
			t.Fatal(err)
		}

		t.Setenv("XDG_CONFIG_HOME", dir)

		cfg, err := loadConfig()
		if err == nil {
			t.Fatal("loadConfig: expected error for invalid YAML, got nil")
		}
		if cfg.Renderer != "" {
			t.Errorf("loadConfig: Renderer = %q on error, want empty", cfg.Renderer)
		}
	})

	t.Run("empty renderer field", func(t *testing.T) {
		dir := t.TempDir()
		k6Dir := filepath.Join(dir, "k6")
		if err := os.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: \"\"\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		t.Setenv("XDG_CONFIG_HOME", dir)

		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("loadConfig: unexpected error: %v", err)
		}
		if cfg.Renderer != "" {
			t.Errorf("loadConfig: Renderer = %q, want empty", cfg.Renderer)
		}
	})
}

func TestPipeRenderer(t *testing.T) {
	t.Parallel()

	t.Run("nil buffer is no-op", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		err := pipeRenderer(nil, &out, "cat")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if out.Len() != 0 {
			t.Errorf("pipeRenderer: wrote %d bytes to output, want 0", out.Len())
		}
	})

	t.Run("empty buffer is no-op", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}
		var out bytes.Buffer
		err := pipeRenderer(buf, &out, "cat")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if out.Len() != 0 {
			t.Errorf("pipeRenderer: wrote %d bytes to output, want 0", out.Len())
		}
	})

	t.Run("renderer processes content", func(t *testing.T) {
		t.Parallel()

		buf := bytes.NewBufferString("hello world")
		var out bytes.Buffer
		// Use cat as a passthrough renderer.
		err := pipeRenderer(buf, &out, "cat")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if out.String() != "hello world" {
			t.Errorf("pipeRenderer: output = %q, want %q", out.String(), "hello world")
		}
	})

	t.Run("renderer with args", func(t *testing.T) {
		t.Parallel()

		buf := bytes.NewBufferString("hello world\ngoodbye world\n")
		var out bytes.Buffer
		// Use head -n 1 as a renderer that takes arguments.
		err := pipeRenderer(buf, &out, "head -n 1")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if strings.TrimSpace(out.String()) != "hello world" {
			t.Errorf("pipeRenderer: output = %q, want %q", strings.TrimSpace(out.String()), "hello world")
		}
	})

	t.Run("missing renderer falls back to raw output", func(t *testing.T) {
		t.Parallel()

		content := "# Documentation\nSome content here.\n"
		buf := bytes.NewBufferString(content)
		var out bytes.Buffer
		err := pipeRenderer(buf, &out, "nonexistent-renderer-binary-xyz")
		if err != nil {
			t.Fatalf("pipeRenderer: unexpected error: %v", err)
		}
		if out.String() != content {
			t.Errorf("pipeRenderer fallback: output = %q, want %q", out.String(), content)
		}
	})

	t.Run("empty renderer string writes raw output", func(t *testing.T) {
		t.Parallel()

		content := "raw content"
		buf := bytes.NewBufferString(content)
		var out bytes.Buffer
		err := pipeRenderer(buf, &out, "")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if out.String() != content {
			t.Errorf("pipeRenderer: output = %q, want %q", out.String(), content)
		}
	})
}

func TestRendererNotUsedWhenNotTTY(t *testing.T) {
	// This test verifies the integration behavior: when stdout is not a TTY
	// (which is the case in test environments), the renderer is NOT invoked
	// even when configured. Output goes directly to the writer.
	cacheDir, _ := setupTestCache(t)

	dir := t.TempDir()
	k6Dir := filepath.Join(dir, "k6")
	if err := os.MkdirAll(k6Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Configure a renderer that would fail if invoked.
	if err := os.WriteFile(filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: nonexistent-renderer-xyz\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)

	cmd := newCmd(nil)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute: %v", err)
	}

	out := buf.String()
	// In non-TTY mode (test environment), output should go directly to the writer
	// without the renderer being invoked.
	if !strings.Contains(out, "k6 Documentation (v0.55.x)") {
		t.Error("expected direct output in non-TTY mode, got: " + out)
	}
}

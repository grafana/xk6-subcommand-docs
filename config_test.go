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

		var stdout, fallback bytes.Buffer
		err := pipeRenderer(nil, &stdout, &fallback, "cat")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if stdout.Len()+fallback.Len() != 0 {
			t.Error("expected no output")
		}
	})

	t.Run("empty buffer is no-op", func(t *testing.T) {
		t.Parallel()

		var stdout, fallback bytes.Buffer
		err := pipeRenderer(&bytes.Buffer{}, &stdout, &fallback, "cat")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if stdout.Len()+fallback.Len() != 0 {
			t.Error("expected no output")
		}
	})

	t.Run("renderer writes to stdout not fallback", func(t *testing.T) {
		t.Parallel()

		buf := bytes.NewBufferString("hello world")
		var stdout, fallback bytes.Buffer
		err := pipeRenderer(buf, &stdout, &fallback, "cat")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if stdout.String() != "hello world" {
			t.Errorf("stdout = %q, want %q", stdout.String(), "hello world")
		}
		if fallback.Len() != 0 {
			t.Errorf("fallback got %d bytes, want 0", fallback.Len())
		}
	})

	t.Run("renderer with args writes to stdout", func(t *testing.T) {
		t.Parallel()

		buf := bytes.NewBufferString("hello world\ngoodbye world\n")
		var stdout, fallback bytes.Buffer
		err := pipeRenderer(buf, &stdout, &fallback, "head -n 1")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if strings.TrimSpace(stdout.String()) != "hello world" {
			t.Errorf("stdout = %q, want %q", strings.TrimSpace(stdout.String()), "hello world")
		}
		if fallback.Len() != 0 {
			t.Errorf("fallback got %d bytes, want 0", fallback.Len())
		}
	})

	t.Run("missing renderer falls back", func(t *testing.T) {
		t.Parallel()

		content := "# Documentation\nSome content here.\n"
		var stdout, fallback bytes.Buffer
		err := pipeRenderer(bytes.NewBufferString(content), &stdout, &fallback, "nonexistent-renderer-binary-xyz")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if fallback.String() != content {
			t.Errorf("fallback = %q, want %q", fallback.String(), content)
		}
		if stdout.Len() != 0 {
			t.Errorf("stdout got %d bytes, want 0", stdout.Len())
		}
	})

	t.Run("failing renderer falls back", func(t *testing.T) {
		t.Parallel()

		content := "# Documentation\nSome content here.\n"
		var stdout, fallback bytes.Buffer
		err := pipeRenderer(bytes.NewBufferString(content), &stdout, &fallback, "false")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if fallback.String() != content {
			t.Errorf("fallback = %q, want %q", fallback.String(), content)
		}
	})

	t.Run("empty renderer string falls back", func(t *testing.T) {
		t.Parallel()

		content := "raw content"
		var stdout, fallback bytes.Buffer
		err := pipeRenderer(bytes.NewBufferString(content), &stdout, &fallback, "")
		if err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		if fallback.String() != content {
			t.Errorf("fallback = %q, want %q", fallback.String(), content)
		}
		if stdout.Len() != 0 {
			t.Errorf("stdout got %d bytes, want 0", stdout.Len())
		}
	})
}

func TestRendererNotUsedWhenNotTTY(t *testing.T) {
	cacheDir, _ := setupTestCache(t)

	dir := t.TempDir()
	k6Dir := filepath.Join(dir, "k6")
	if err := os.MkdirAll(k6Dir, 0o755); err != nil {
		t.Fatal(err)
	}
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
	if !strings.Contains(out, "k6 Documentation (v0.55.x)") {
		t.Error("expected direct output in non-TTY mode, got: " + out)
	}
}

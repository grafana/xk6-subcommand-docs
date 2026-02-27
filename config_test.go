package docs

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"go.k6.io/k6/lib/fsext"
)

func TestConfigDir(t *testing.T) {
	t.Parallel()

	t.Run("XDG_CONFIG_HOME has priority", func(t *testing.T) {
		t.Parallel()

		env := map[string]string{
			"XDG_CONFIG_HOME": "/xdg",
			"HOME":            "/home/fallback",
			"USERPROFILE":     "/users/fallback",
		}
		dir, err := configDir(env)
		if err != nil {
			t.Fatalf("configDir: unexpected error: %v", err)
		}
		want := filepath.Join("/xdg", "k6")
		if dir != want {
			t.Errorf("configDir = %q, want %q", dir, want)
		}
	})

	t.Run("HOME preferred over USERPROFILE", func(t *testing.T) {
		t.Parallel()

		env := map[string]string{
			"HOME":        "/home/test",
			"USERPROFILE": "/users/test",
		}
		dir, err := configDir(env)
		if err != nil {
			t.Fatalf("configDir: unexpected error: %v", err)
		}
		want := filepath.Join("/home/test", ".config", "k6")
		if dir != want {
			t.Errorf("configDir = %q, want %q", dir, want)
		}
	})

	t.Run("USERPROFILE fallback when HOME is unset", func(t *testing.T) {
		t.Parallel()

		env := map[string]string{"USERPROFILE": "/users/test"}
		dir, err := configDir(env)
		if err != nil {
			t.Fatalf("configDir: unexpected error: %v", err)
		}
		want := filepath.Join("/users/test", ".config", "k6")
		if dir != want {
			t.Errorf("configDir = %q, want %q", dir, want)
		}
	})

	t.Run("error when neither HOME nor USERPROFILE is set", func(t *testing.T) {
		t.Parallel()

		env := map[string]string{"XDG_CONFIG_HOME": ""}
		_, err := configDir(env)
		if err == nil {
			t.Fatal("configDir: expected error when neither HOME nor USERPROFILE is set")
		}
	})
}

func TestCacheDirUSERPROFILE(t *testing.T) {
	t.Parallel()

	env := map[string]string{"USERPROFILE": "/users/test"}
	dir, err := CacheDir(env, "v1.0.0")
	if err != nil {
		t.Fatalf("CacheDir: unexpected error: %v", err)
	}
	want := filepath.Join("/users/test", ".local", "share", "k6", "docs", "v1.0.0")
	if dir != want {
		t.Errorf("CacheDir = %q, want %q", dir, want)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid config", func(t *testing.T) {
		t.Parallel()
		afs := fsext.NewMemMapFs()
		dir := "/tmp/config"
		env := map[string]string{"XDG_CONFIG_HOME": dir}

		k6Dir := filepath.Join(dir, "k6")
		if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: glow -p 200\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := loadConfig(afs, env)
		if err != nil {
			t.Fatalf("loadConfig: unexpected error: %v", err)
		}
		if cfg.Renderer != "glow -p 200" {
			t.Errorf("loadConfig: Renderer = %q, want %q", cfg.Renderer, "glow -p 200")
		}
	})

	t.Run("missing file returns empty config", func(t *testing.T) {
		t.Parallel()
		afs := fsext.NewMemMapFs()
		dir := "/tmp/config-missing"
		env := map[string]string{"XDG_CONFIG_HOME": dir}

		cfg, err := loadConfig(afs, env)
		if err != nil {
			t.Fatalf("loadConfig: unexpected error: %v", err)
		}
		if cfg.Renderer != "" {
			t.Errorf("loadConfig: Renderer = %q, want empty", cfg.Renderer)
		}
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		t.Parallel()
		afs := fsext.NewMemMapFs()
		dir := "/tmp/config-invalid"
		env := map[string]string{"XDG_CONFIG_HOME": dir}

		k6Dir := filepath.Join(dir, "k6")
		if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte(":\n  :\n    : [invalid"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := loadConfig(afs, env)
		if err == nil {
			t.Fatal("loadConfig: expected error for invalid YAML, got nil")
		}
		if cfg.Renderer != "" {
			t.Errorf("loadConfig: Renderer = %q on error, want empty", cfg.Renderer)
		}
	})

	t.Run("config via HOME fallback", func(t *testing.T) {
		t.Parallel()
		afs := fsext.NewMemMapFs()
		home := "/home/testuser"
		env := map[string]string{"HOME": home}

		k6Dir := filepath.Join(home, ".config", "k6")
		if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: glow -w 200\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := loadConfig(afs, env)
		if err != nil {
			t.Fatalf("loadConfig: unexpected error: %v", err)
		}
		if cfg.Renderer != "glow -w 200" {
			t.Errorf("loadConfig: Renderer = %q, want %q", cfg.Renderer, "glow -w 200")
		}
	})

	t.Run("empty renderer field", func(t *testing.T) {
		t.Parallel()
		afs := fsext.NewMemMapFs()
		dir := "/tmp/config-empty"
		env := map[string]string{"XDG_CONFIG_HOME": dir}

		k6Dir := filepath.Join(dir, "k6")
		if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: \"\"\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := loadConfig(afs, env)
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

	mustPipe := func(t *testing.T, buf *bytes.Buffer, renderer string) (stdout, fallback string) {
		t.Helper()
		var stdoutBuf, fallbackBuf bytes.Buffer
		if err := pipeRenderer(context.Background(), buf, &stdoutBuf, &fallbackBuf, io.Discard, renderer); err != nil {
			t.Fatalf("pipeRenderer: %v", err)
		}
		return stdoutBuf.String(), fallbackBuf.String()
	}

	t.Run("nil buffer is no-op", func(t *testing.T) {
		t.Parallel()
		stdout, fallback := mustPipe(t, nil, "cat")
		if len(stdout)+len(fallback) != 0 {
			t.Error("expected no output")
		}
	})

	t.Run("empty buffer is no-op", func(t *testing.T) {
		t.Parallel()
		stdout, fallback := mustPipe(t, &bytes.Buffer{}, "cat")
		if len(stdout)+len(fallback) != 0 {
			t.Error("expected no output")
		}
	})

	t.Run("renderer writes to stdout not fallback", func(t *testing.T) {
		t.Parallel()
		stdout, fallback := mustPipe(t, bytes.NewBufferString("hello world"), "cat")
		if stdout != "hello world" {
			t.Errorf("stdout = %q, want %q", stdout, "hello world")
		}
		if fallback != "" {
			t.Errorf("fallback = %q, want empty", fallback)
		}
	})

	t.Run("renderer with args writes to stdout", func(t *testing.T) {
		t.Parallel()
		stdout, fallback := mustPipe(t, bytes.NewBufferString("hello world\ngoodbye world\n"), "head -n 1")
		if strings.TrimSpace(stdout) != "hello world" {
			t.Errorf("stdout = %q, want %q", strings.TrimSpace(stdout), "hello world")
		}
		if fallback != "" {
			t.Errorf("fallback = %q, want empty", fallback)
		}
	})

	t.Run("missing renderer falls back", func(t *testing.T) {
		t.Parallel()
		content := "# Documentation\nSome content here.\n"
		stdout, fallback := mustPipe(t, bytes.NewBufferString(content), "nonexistent-renderer-binary-xyz")
		if fallback != content {
			t.Errorf("fallback = %q, want %q", fallback, content)
		}
		if stdout != "" {
			t.Errorf("stdout = %q, want empty", stdout)
		}
	})

	t.Run("failing renderer falls back", func(t *testing.T) {
		t.Parallel()
		content := "# Documentation\nSome content here.\n"
		_, fallback := mustPipe(t, bytes.NewBufferString(content), "false")
		if fallback != content {
			t.Errorf("fallback = %q, want %q", fallback, content)
		}
	})

	t.Run("empty renderer string falls back", func(t *testing.T) {
		t.Parallel()
		content := "raw content"
		stdout, fallback := mustPipe(t, bytes.NewBufferString(content), "")
		if fallback != content {
			t.Errorf("fallback = %q, want %q", fallback, content)
		}
		if stdout != "" {
			t.Errorf("stdout = %q, want empty", stdout)
		}
	})
}

func TestSearchUsesRenderer(t *testing.T) {
	t.Parallel()
	afs, cacheDir, _ := setupTestCache(t)

	gs := newTestGlobalState(t, afs)
	env := gs.Env
	env["XDG_CONFIG_HOME"] = "/tmp/search-renderer-config"

	k6Dir := filepath.Join(env["XDG_CONFIG_HOME"], "k6")
	if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Configure renderer as "cat -n" which adds line numbers — if the
	// renderer is used, output will contain line-numbered text.
	if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: cat -n\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Force TTY so the renderer path is triggered.
	gs.Stdout.IsTTY = true

	// pipeRenderer writes rendered output to gs.Stdout, so capture it.
	var stdoutBuf bytes.Buffer
	gs.Stdout.Writer = &stdoutBuf

	cmd := newCmd(gs)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "search", "http"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute: %v", err)
	}

	out := stdoutBuf.String()
	// Search should find results for "http".
	if !strings.Contains(out, "http") {
		t.Errorf("expected search results for 'http', got: %s", out)
	}
	// If the renderer is used, "cat -n" prepends line numbers like "     1\t".
	if !strings.Contains(out, "\t") {
		t.Error("expected renderer (cat -n) to be used for search output, but output has no tab characters from line numbering")
	}
}

func TestRendererUsedForTopicWhenTTY(t *testing.T) {
	t.Parallel()
	afs, cacheDir, _ := setupTestCache(t)

	gs := newTestGlobalState(t, afs)
	env := gs.Env
	// Use HOME (not XDG_CONFIG_HOME) — this is the path real users hit.
	env["HOME"] = "/home/testuser"

	k6Dir := filepath.Join(env["HOME"], ".config", "k6")
	if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: cat -n\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	gs.Stdout.IsTTY = true

	var stdoutBuf bytes.Buffer
	gs.Stdout.Writer = &stdoutBuf

	cmd := newCmd(gs)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "http", "get"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute: %v", err)
	}

	out := stdoutBuf.String()
	if !strings.Contains(out, "get") {
		t.Errorf("expected topic content for 'http get', got: %s", out)
	}
	if !strings.Contains(out, "\t") {
		t.Error("expected renderer (cat -n) to be used for topic output, but output has no tab characters from line numbering")
	}
}

func TestRendererPreservesOutput(t *testing.T) {
	t.Parallel()

	// Run without renderer to get the expected output.
	afs, cacheDir, _ := setupTestCache(t)
	gs := newTestGlobalState(t, afs)

	cmd := newCmd(gs)
	var expectedBuf bytes.Buffer
	cmd.SetOut(&expectedBuf)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--cache-dir", cacheDir, "--version", "v0.55.x", "http", "get"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute (no renderer): %v", err)
	}
	expected := expectedBuf.String()

	// Run with renderer "cat" and TTY on — output should be identical.
	afs2, cacheDir2, _ := setupTestCache(t)
	gs2 := newTestGlobalState(t, afs2)
	gs2.Env["XDG_CONFIG_HOME"] = "/tmp/renderer-preserve-config"
	gs2.Stdout.IsTTY = true

	k6Dir := filepath.Join(gs2.Env["XDG_CONFIG_HOME"], "k6")
	if err := afs2.MkdirAll(k6Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := fsext.WriteFile(afs2, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: cat\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var gotBuf bytes.Buffer
	gs2.Stdout.Writer = &gotBuf

	cmd2 := newCmd(gs2)
	cmd2.SetErr(io.Discard)
	cmd2.SetArgs([]string{"--cache-dir", cacheDir2, "--version", "v0.55.x", "http", "get"})
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("cmd.Execute (with renderer): %v", err)
	}
	got := gotBuf.String()

	if got != expected {
		t.Errorf("renderer altered output:\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestRendererNotUsedWhenNotTTY(t *testing.T) {
	t.Parallel()
	afs, cacheDir, _ := setupTestCache(t)

	gs := newTestGlobalState(t, afs)
	env := gs.Env
	env["XDG_CONFIG_HOME"] = "/tmp/renderertest-config"

	k6Dir := filepath.Join(env["XDG_CONFIG_HOME"], "k6")
	if err := afs.MkdirAll(k6Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := fsext.WriteFile(afs, filepath.Join(k6Dir, "docs.yaml"), []byte("renderer: nonexistent-renderer-xyz\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := newCmd(gs)
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

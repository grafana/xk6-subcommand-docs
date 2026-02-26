package docs

import (
	"strings"
	"testing"
)

func TestStripFrontmatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no frontmatter",
			input: "# Hello\n\nWorld",
			want:  "# Hello\n\nWorld",
		},
		{
			name:  "with frontmatter",
			input: "---\ntitle: 'Test'\nweight: 10\n---\n\n# Hello",
			want:  "\n# Hello",
		},
		{
			name:  "empty content after frontmatter",
			input: "---\ntitle: 'X'\n---\n",
			want:  "",
		},
		{
			name:  "only frontmatter",
			input: "---\ntitle: 'X'\n---",
			want:  "",
		},
		{
			name:  "unclosed frontmatter treated as no frontmatter",
			input: "---\ntitle: 'X'\nno closing",
			want:  "---\ntitle: 'X'\nno closing",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "frontmatter not at start is left alone",
			input: "hello\n---\ntitle: 'X'\n---\nworld",
			want:  "hello\n---\ntitle: 'X'\n---\nworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := StripFrontmatter(tt.input)
			if got != tt.want {
				t.Errorf("StripFrontmatter() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTransform_ResolveShared(t *testing.T) {
	t.Parallel()

	shared := map[string]string{
		"javascript-api/k6-http.md": "---\ntitle: shared\n---\n\nThe k6/http module handles HTTP.",
		"preview-feature.md":        "---\ntitle: preview\n---\n\nThis is a preview feature.",
	}

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "resolve shared shortcode",
			content: `{{< docs/shared source="k6" lookup="javascript-api/k6-http.md" version="<K6_VERSION>" >}}`,
			want:    "\nThe k6/http module handles HTTP.",
		},
		{
			name:    "resolve shared shortcode with different spacing",
			content: `{{<  docs/shared  source="k6"  lookup="preview-feature.md"  version="<K6_VERSION>"  >}}`,
			want:    "\nThis is a preview feature.",
		},
		{
			name:    "missing shared content is removed",
			content: `{{< docs/shared source="k6" lookup="nonexistent.md" version="<K6_VERSION>" >}}`,
			want:    "",
		},
		{
			name:    "shared content inline with surrounding text",
			content: "Before\n\n{{< docs/shared source=\"k6\" lookup=\"javascript-api/k6-http.md\" version=\"<K6_VERSION>\" >}}\n\nAfter",
			want:    "Before\n\n\nThe k6/http module handles HTTP.\n\nAfter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.5.x", shared)
			// Only check the shared resolution part — other steps may also run.
			// We check that the shared content is present and shortcode is gone.
			if !strings.Contains(got, "k6/http module handles HTTP") && tt.name == "resolve shared shortcode" {
				t.Errorf("expected shared content to be inlined, got: %q", got)
			}
			if strings.Contains(got, "docs/shared") {
				t.Errorf("shortcode should have been removed, got: %q", got)
			}
		})
	}
}

func TestTransform_StripCodeTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "code tags around code block",
			content: "{{< code >}}\n\n```javascript\nconsole.log('hi');\n```\n\n{{< /code >}}",
			want:    "\n\n```javascript\nconsole.log('hi');\n```\n\n",
		},
		{
			name:    "multiple code tag pairs",
			content: "{{< code >}}\n```js\na();\n```\n{{< /code >}}\n\n{{< code >}}\n```js\nb();\n```\n{{< /code >}}",
			want:    "\n```js\na();\n```\n\n```js\nb();\n```\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.0.0", nil)
			if got != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestTransform_Admonition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "note admonition",
			content: `{{< admonition type="note" >}}

This is a note.

{{< /admonition >}}`,
			want: "> **Note:** This is a note.\n",
		},
		{
			name: "warning admonition",
			content: `{{< admonition type="warning" >}}

Be careful with this.

{{< /admonition >}}`,
			want: "> **Warning:** Be careful with this.\n",
		},
		{
			name: "caution admonition",
			content: `{{< admonition type="caution" >}}

Caution text here.

{{< /admonition >}}`,
			want: "> **Caution:** Caution text here.\n",
		},
		{
			name: "multiline admonition",
			content: `{{< admonition type="note" >}}

Line one.
Line two.
Line three.

{{< /admonition >}}`,
			want: "> **Note:** Line one.\n> Line two.\n> Line three.\n",
		},
		{
			name: "admonition with surrounding content",
			content: "Before\n\n{{< admonition type=\"note\" >}}\n\nImportant thing.\n\n{{< /admonition >}}\n\nAfter",
			want:    "Before\n\n> **Note:** Important thing.\n\nAfter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.0.0", nil)
			if got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestTransform_StripSection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "simple section",
			content: "Before\n\n{{< section >}}\n\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "section with attributes",
			content: "Before\n\n{{< section depth=2 >}}\n\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "section with menuTitle",
			content: `Before` + "\n\n" + `{{< section menuTitle="true">}}` + "\n\nAfter",
			want:    "Before\n\nAfter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.0.0", nil)
			if got != tt.want {
				t.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

func TestTransform_StripRemainingShortcodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "youtube shortcode",
			content: "Before\n\n{{< youtube id=\"1mtYVDA2_iQ\" >}}\n\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "card-grid shortcode",
			content: "Before\n\n{{< card-grid key=\"cards\" type=\"simple\" >}}\n\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "docs/hero-simple shortcode",
			content: "Before\n\n{{< docs/hero-simple key=\"hero\" >}}\n\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "collapse opening and closing",
			content: "Before\n\n{{< collapse title=\"example\" >}}\n\nContent\n\n{{< /collapse >}}\n\nAfter",
			want:    "Before\n\nContent\n\nAfter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.0.0", nil)
			if got != tt.want {
				t.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

func TestTransform_ReplaceVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		version string
		want    string
	}{
		{
			name:    "replace version in bare URL",
			content: "Visit https://grafana.com/docs/k6/<K6_VERSION>/extensions/explore for extensions.",
			version: "v1.5.x",
			want:    "Visit https://grafana.com/docs/k6/v1.5.x/extensions/explore for extensions.",
		},
		{
			name:    "replace version in external link (kept as link)",
			content: "[extensions](https://grafana.com/docs/k6/<K6_VERSION>/extensions/explore)",
			version: "v1.5.x",
			want:    "[extensions](https://grafana.com/docs/k6/v1.5.x/extensions/explore)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, tt.version, nil)
			if got != tt.want {
				t.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

func TestTransform_ConvertInternalLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "internal JS API link becomes plain text",
			content: "[batch( requests )](https://grafana.com/docs/k6/v1.5.x/javascript-api/k6-http/batch)",
			want:    "batch( requests )",
		},
		{
			name:    "internal using-k6 link becomes plain text",
			content: "[thresholds](https://grafana.com/docs/k6/v1.5.x/using-k6/thresholds)",
			want:    "thresholds",
		},
		{
			name:    "internal link with anchor becomes plain text",
			content: "[URL Grouping](https://grafana.com/docs/k6/v1.5.x/using-k6/http-requests#url-grouping)",
			want:    "URL Grouping",
		},
		{
			name:    "internal link with trailing slash becomes plain text",
			content: "[scenarios](https://grafana.com/docs/k6/v1.5.x/using-k6/scenarios/)",
			want:    "scenarios",
		},
		{
			name:    "external link to extensions keeps URL",
			content: "[Build a k6 binary](https://grafana.com/docs/k6/v1.5.x/extensions/build-k6-binary-using-go)",
			want:    "[Build a k6 binary](https://grafana.com/docs/k6/v1.5.x/extensions/build-k6-binary-using-go)",
		},
		{
			name:    "external link to get-started keeps URL",
			content: "[Install k6](https://grafana.com/docs/k6/v1.5.x/get-started/installation/)",
			want:    "[Install k6](https://grafana.com/docs/k6/v1.5.x/get-started/installation/)",
		},
		{
			name:    "external link to set-up keeps URL",
			content: "[Set up](https://grafana.com/docs/k6/v1.5.x/set-up/something)",
			want:    "[Set up](https://grafana.com/docs/k6/v1.5.x/set-up/something)",
		},
		{
			name:    "multiple links in one line",
			content: "See [metrics](https://grafana.com/docs/k6/v1.5.x/using-k6/metrics) and [checks](https://grafana.com/docs/k6/v1.5.x/using-k6/checks).",
			want:    "See metrics and checks.",
		},
		{
			name:    "link text with brackets (optional params)",
			content: "[get( url, [params] )](https://grafana.com/docs/k6/v1.5.x/javascript-api/k6-http/get)",
			want:    "get( url, [params] )",
		},
		{
			name:    "link text with nested brackets",
			content: "[check(selector[, options])](https://grafana.com/docs/k6/v1.5.x/javascript-api/k6-browser/page/check/)",
			want:    "check(selector[, options])",
		},
		{
			name:    "non-grafana link left alone",
			content: "[example](https://example.com/something)",
			want:    "[example](https://example.com/something)",
		},
		{
			name:    "all included categories become plain text",
			content: "[a](https://grafana.com/docs/k6/v1.5.x/javascript-api/foo) [b](https://grafana.com/docs/k6/v1.5.x/using-k6/bar) [c](https://grafana.com/docs/k6/v1.5.x/using-k6-browser/baz) [d](https://grafana.com/docs/k6/v1.5.x/testing-guides/qux) [e](https://grafana.com/docs/k6/v1.5.x/examples/quux) [f](https://grafana.com/docs/k6/v1.5.x/results-output/corge) [g](https://grafana.com/docs/k6/v1.5.x/reference/grault)",
			want:    "a b c d e f g",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.5.x", nil)
			if got != tt.want {
				t.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

func TestTransform_StripHTMLComments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "single line comment",
			content: "Before\n<!-- md-k6:skip -->\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "multiline comment",
			content: "Before\n<!--\nThis is\na multiline comment\n-->\nAfter",
			want:    "Before\n\nAfter",
		},
		{
			name:    "inline comment",
			content: "Hello <!-- hidden --> world",
			want:    "Hello  world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.0.0", nil)
			if got != tt.want {
				t.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

func TestTransform_StripFrontmatter(t *testing.T) {
	t.Parallel()

	content := "---\ntitle: 'Checks'\ndescription: 'Some description.'\nweight: 400\n---\n\n# Checks\n\nContent here."
	want := "\n# Checks\n\nContent here."

	got := Transform(content, "v1.0.0", nil)
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestTransform_NormalizeWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "three newlines collapsed to two",
			content: "A\n\n\nB",
			want:    "A\n\nB",
		},
		{
			name:    "four newlines collapsed to two",
			content: "A\n\n\n\nB",
			want:    "A\n\nB",
		},
		{
			name:    "many newlines collapsed",
			content: "A\n\n\n\n\n\n\nB",
			want:    "A\n\nB",
		},
		{
			name:    "two newlines kept as is",
			content: "A\n\nB",
			want:    "A\n\nB",
		},
		{
			name:    "single newline kept",
			content: "A\nB",
			want:    "A\nB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Transform(tt.content, "v1.0.0", nil)
			if got != tt.want {
				t.Errorf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

func TestTransform_EndToEnd_ChecksMd(t *testing.T) {
	t.Parallel()

	// Simulate the checks.md file from k6-docs (simplified but realistic).
	input := `---
title: 'Checks'
description: 'Checks are like asserts but differ in that they do not halt the execution.'
weight: 400
---

# Checks

Checks validate boolean conditions in your test.

Each check creates a [rate metric](https://grafana.com/docs/k6/<K6_VERSION>/using-k6/metrics).

## Check for HTTP response code

{{< code >}}

` + "```javascript" + `
import { check } from 'k6';
import http from 'k6/http';

export default function () {
  const res = http.get('http://test.k6.io/');
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
}
` + "```" + `

{{< /code >}}

{{< admonition type="note" >}}

When a check fails, the script will continue executing successfully and will not return a 'failed' exit status.
If you need the whole test to fail based on the results of a check, you have to combine checks with thresholds.

{{< /admonition >}}

<!-- internal tracking comment -->

## Read more

- [Check API](https://grafana.com/docs/k6/<K6_VERSION>/javascript-api/k6/check)`

	got := Transform(input, "v1.5.x", nil)

	// Should not contain frontmatter.
	if strings.Contains(got, "title: 'Checks'") {
		t.Error("frontmatter should be stripped")
	}

	// Should not contain code shortcodes.
	if strings.Contains(got, "{{< code >}}") || strings.Contains(got, "{{< /code >}}") {
		t.Error("code shortcodes should be stripped")
	}

	// Should contain the code block content.
	if !strings.Contains(got, "import { check } from 'k6';") {
		t.Error("code block content should be preserved")
	}

	// Admonition should be converted to blockquote.
	if !strings.Contains(got, "> **Note:**") {
		t.Error("admonition should be converted to blockquote")
	}
	if strings.Contains(got, "{{< admonition") || strings.Contains(got, "{{< /admonition") {
		t.Error("admonition shortcodes should be stripped")
	}

	// Version should be replaced.
	if strings.Contains(got, "<K6_VERSION>") {
		t.Error("version placeholder should be replaced")
	}

	// Internal docs links should be plain text (no URL).
	if strings.Contains(got, "grafana.com/docs/k6") {
		t.Error("internal docs links should be converted to plain text")
	}
	if !strings.Contains(got, "rate metric") {
		t.Error("link text should be preserved")
	}

	// HTML comments should be stripped.
	if strings.Contains(got, "<!-- internal") {
		t.Error("HTML comments should be stripped")
	}

	// No triple+ newlines.
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("should not have 3+ consecutive newlines, got:\n%s", got)
	}
}

func TestTransform_EndToEnd_SharedContent(t *testing.T) {
	t.Parallel()

	shared := map[string]string{
		"javascript-api/k6-http.md": "---\ntitle: 'k6/http shared'\n---\n\nThe k6/http module contains functionality for performing HTTP transactions.\n\n| Method | Description |\n|--------|-------------|\n| get    | GET request |",
	}

	input := `---
title: 'k6/http'
description: 'The k6/http module contains functionality for performing HTTP transactions.'
weight: 09
---

# k6/http

{{< docs/shared source="k6" lookup="javascript-api/k6-http.md" version="<K6_VERSION>" >}}`

	got := Transform(input, "v1.5.x", shared)

	if strings.Contains(got, "title: 'k6/http'") {
		t.Error("frontmatter should be stripped")
	}

	if strings.Contains(got, "docs/shared") {
		t.Error("shared shortcode should be resolved")
	}

	if !strings.Contains(got, "The k6/http module contains functionality") {
		t.Error("shared content should be inlined")
	}

	if !strings.Contains(got, "| get    | GET request |") {
		t.Error("shared content table should be inlined")
	}

	// The shared content's own frontmatter should be stripped.
	if strings.Contains(got, "title: 'k6/http shared'") {
		t.Error("shared content frontmatter should be stripped")
	}
}

func TestTransform_EndToEnd_ScenariosMd(t *testing.T) {
	t.Parallel()

	input := `---
title: Scenarios
description: 'Scenarios configure how VUs and iteration schedules.'
weight: 1500
---

# Scenarios

Scenarios configure how VUs and iteration schedules in granular detail.

{{< code >}}

<!-- md-k6:skip -->

` + "```javascript" + `
export const options = {
  scenarios: {
    example_scenario: {
      executor: 'shared-iterations',
    },
  },
};
` + "```" + `

{{< /code >}}

## Scenario executors

See [open vs closed](https://grafana.com/docs/k6/<K6_VERSION>/using-k6/scenarios/concepts/open-vs-closed).

{{< section >}}

{{< collapse title="full output" >}}

` + "```bash" + `
running (00m12.8s), 00/20 VUs
` + "```" + `

{{< /collapse >}}`

	got := Transform(input, "v1.5.x", nil)

	// Frontmatter gone.
	if strings.Contains(got, "title: Scenarios") {
		t.Error("frontmatter should be stripped")
	}

	// Code tags gone, content preserved.
	if strings.Contains(got, "{{< code >}}") {
		t.Error("code shortcodes should be stripped")
	}
	if !strings.Contains(got, "executor: 'shared-iterations'") {
		t.Error("code block content should be preserved")
	}

	// HTML comment gone.
	if strings.Contains(got, "<!-- md-k6:skip -->") {
		t.Error("HTML comments should be stripped")
	}

	// Section gone.
	if strings.Contains(got, "{{< section >}}") {
		t.Error("section shortcodes should be stripped")
	}

	// Version replaced and internal link converted to plain text.
	if strings.Contains(got, "<K6_VERSION>") {
		t.Error("version placeholder should be replaced")
	}
	if strings.Contains(got, "grafana.com/docs/k6") {
		t.Error("internal docs links should be converted to plain text")
	}
	if !strings.Contains(got, "open vs closed") {
		t.Error("link text should be preserved")
	}

	// Collapse tags stripped, content preserved.
	if strings.Contains(got, "{{< collapse") || strings.Contains(got, "{{< /collapse") {
		t.Error("collapse shortcodes should be stripped")
	}
	if !strings.Contains(got, "running (00m12.8s)") {
		t.Error("collapse content should be preserved")
	}

	// No triple newlines.
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("should not have 3+ consecutive newlines, got:\n%s", got)
	}
}

func TestTransform_EmptyInput(t *testing.T) {
	t.Parallel()

	got := Transform("", "v1.0.0", nil)
	if got != "" {
		t.Errorf("expected empty string, got: %q", got)
	}
}

func TestTransform_NoOp(t *testing.T) {
	t.Parallel()

	input := "# Simple markdown\n\nJust text, nothing special."
	got := Transform(input, "v1.0.0", nil)
	if got != input {
		t.Errorf("got: %q, want: %q", got, input)
	}
}

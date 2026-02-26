package docs

import (
	"regexp"
	"strings"
)

var (
	reShared       = regexp.MustCompile(`\{\{<\s*docs/shared\s+source="k6"\s+lookup="([^"]+)".*?>\}\}`)
	reCodeTag      = regexp.MustCompile(`\{\{<\s*/?\s*code\s*>\}\}`)
	reAdmonition   = regexp.MustCompile(`(?s)\{\{<\s*admonition\s+type="([^"]+)"\s*>\}\}\s*\n(.*?)\n\s*\{\{<\s*/admonition\s*>\}\}`)
	reSection      = regexp.MustCompile(`\{\{<\s*/?\s*section\b[^>]*>\}\}`)
	reAnyShortcode = regexp.MustCompile(`\{\{<\s*/?\s*[^>]+>\}\}`)
	reHTMLComment  = regexp.MustCompile(`<!--[\s\S]*?-->`)
	reExtraNewline = regexp.MustCompile(`\n{3,}`)
)

// Transform applies Hugo shortcode resolution and markdown cleanup to content.
// The pipeline runs in a fixed order:
//  1. Resolve docs/shared shortcodes
//  2. Strip code tags
//  3. Convert admonitions to blockquotes
//  4. Strip section tags
//  5. Strip remaining shortcodes
//  6. Replace <K6_VERSION> with version
//  7. Strip HTML comments
//  8. Strip YAML frontmatter
//  9. Normalize whitespace
func Transform(content, version string, sharedContent map[string]string) string {
	if content == "" {
		return ""
	}

	s := content

	// 1. Resolve shared shortcodes.
	s = reShared.ReplaceAllStringFunc(s, func(match string) string {
		m := reShared.FindStringSubmatch(match)
		if m == nil || sharedContent == nil {
			return ""
		}
		raw, ok := sharedContent[m[1]]
		if !ok {
			return ""
		}
		return StripFrontmatter(raw)
	})

	// 2. Strip code tags (keep content between them).
	s = reCodeTag.ReplaceAllString(s, "")

	// 3. Convert admonitions to blockquotes.
	s = reAdmonition.ReplaceAllStringFunc(s, func(match string) string {
		m := reAdmonition.FindStringSubmatch(match)
		if m == nil {
			return match
		}
		title := strings.ToUpper(m[1][:1]) + m[1][1:]
		body := strings.TrimSpace(m[2])

		lines := strings.Split(body, "\n")
		var sb strings.Builder
		first := true
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if first {
				sb.WriteString("> **" + title + ":** " + line + "\n")
				first = false
			} else {
				sb.WriteString("> " + line + "\n")
			}
		}
		return sb.String()
	})

	// 4. Strip section tags.
	s = reSection.ReplaceAllString(s, "")

	// 5. Strip remaining shortcodes.
	s = reAnyShortcode.ReplaceAllString(s, "")

	// 6. Replace version placeholder.
	s = strings.ReplaceAll(s, "<K6_VERSION>", version)

	// 7. Strip HTML comments.
	s = reHTMLComment.ReplaceAllString(s, "")

	// 8. Strip YAML frontmatter.
	s = StripFrontmatter(s)

	// 9. Normalize whitespace: collapse 3+ consecutive newlines to 2.
	s = reExtraNewline.ReplaceAllString(s, "\n\n")

	return s
}

// StripFrontmatter removes YAML frontmatter (delimited by "---") from the
// start of content. If the content doesn't start with "---\n" or the closing
// delimiter is missing, it returns the content unchanged.
func StripFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---\n") {
		return content
	}

	end := strings.Index(content[4:], "\n---")
	if end == -1 {
		return content
	}

	// Skip past the closing "\n---" (4 bytes).
	cutAt := 4 + end + 4
	// Also consume the newline right after the closing "---" if present.
	if cutAt < len(content) && content[cutAt] == '\n' {
		cutAt++
	}
	return content[cutAt:]
}

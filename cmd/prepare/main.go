// Command prepare processes the k6-docs repository into a doc bundle
// suitable for embedding. It walks the documentation tree, transforms
// Hugo shortcodes into clean markdown, and produces:
//   - markdown/ — transformed .md files
//   - sections.json — structured index of all sections
//   - best_practices.md — a comprehensive best practices guide
package main

import (
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
	bySlug := make(map[string]*docs.Section, len(sections))
	for i := range sections {
		bySlug[sections[i].Slug] = &sections[i]
	}

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

const bestPracticesContent = `# k6 Best Practices

A comprehensive guide to writing effective, maintainable, and performant k6 load tests.

## Test Structure and Organization

### Use the k6 Test Lifecycle

k6 has four distinct lifecycle stages. Use them intentionally:

` + "```javascript" + `
// 1. init — runs once per VU, used to set up test data and imports
import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'https://test-api.example.com';

export const options = {
  stages: [
    { duration: '1m', target: 20 },
    { duration: '3m', target: 20 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    checks: ['rate>0.99'],
  },
};

// 2. setup — runs once before test, used for auth tokens, seed data, etc.
export function setup() {
  const loginRes = http.post(BASE_URL + '/auth/login', JSON.stringify({
    username: 'testuser',
    password: 'testpass',
  }), { headers: { 'Content-Type': 'application/json' } });

  const token = loginRes.json('token');
  return { token };
}

// 3. default function — runs repeatedly for each VU iteration
export default function (data) {
  const params = {
    headers: { Authorization: 'Bearer ' + data.token },
  };

  const res = http.get(BASE_URL + '/api/items', params);
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response has items': (r) => r.json('items').length > 0,
  });

  sleep(1);
}

// 4. teardown — runs once after test, used for cleanup
export function teardown(data) {
  http.post(BASE_URL + '/auth/logout', null, {
    headers: { Authorization: 'Bearer ' + data.token },
  });
}
` + "```" + `

### Group Related Requests

Use ` + "`group()`" + ` to organize related requests into logical transaction blocks:

` + "```javascript" + `
import { group, check } from 'k6';
import http from 'k6/http';

export default function () {
  group('user registration flow', function () {
    const signupRes = http.post('https://api.example.com/signup', JSON.stringify({
      email: 'user@example.com',
      password: 'securepass',
    }), { headers: { 'Content-Type': 'application/json' } });
    check(signupRes, { 'signup status 201': (r) => r.status === 201 });

    const verifyRes = http.get('https://api.example.com/verify?token=abc');
    check(verifyRes, { 'verify status 200': (r) => r.status === 200 });
  });

  group('user login flow', function () {
    const loginRes = http.post('https://api.example.com/login', JSON.stringify({
      email: 'user@example.com',
      password: 'securepass',
    }), { headers: { 'Content-Type': 'application/json' } });
    check(loginRes, { 'login status 200': (r) => r.status === 200 });
  });
}
` + "```" + `

## Performance and Resource Management

### Add Think Time Between Requests

Real users do not fire requests instantly. Use ` + "`sleep()`" + ` to simulate realistic pacing:

` + "```javascript" + `
import { sleep } from 'k6';
import http from 'k6/http';

export default function () {
  http.get('https://test.k6.io/');
  sleep(Math.random() * 3 + 1); // 1-4 seconds of random think time
}
` + "```" + `

### Use Checks Instead of console.log

Avoid ` + "`console.log()`" + ` for validations — it does not integrate with k6 metrics. Use ` + "`check()`" + ` instead:

` + "```javascript" + `
import { check } from 'k6';
import http from 'k6/http';

export default function () {
  const res = http.get('https://test.k6.io/');

  // Bad: does not track pass/fail metrics
  // console.log('Status:', res.status);

  // Good: integrates with thresholds and summary
  check(res, {
    'status is 200': (r) => r.status === 200,
    'body is not empty': (r) => r.body.length > 0,
  });
}
` + "```" + `

### Set Thresholds for Pass/Fail Criteria

Define clear thresholds so CI pipelines can gate on performance:

` + "```javascript" + `
export const options = {
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
    checks: ['rate>0.99'],
    'http_req_duration{name:login}': ['p(95)<800'],
  },
};
` + "```" + `

## Error Handling Patterns

### Validate Responses Thoroughly

Check more than just the status code:

` + "```javascript" + `
import { check } from 'k6';
import http from 'k6/http';

export default function () {
  const res = http.get('https://api.example.com/users/1');

  check(res, {
    'status is 200': (r) => r.status === 200,
    'content-type is json': (r) => r.headers['Content-Type'].includes('application/json'),
    'body has user id': (r) => r.json('id') !== undefined,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
` + "```" + `

### Handle Expected Errors Gracefully

Not all non-200 responses are failures. Design checks for your actual expectations:

` + "```javascript" + `
import { check } from 'k6';
import http from 'k6/http';

export default function () {
  // This endpoint might return 404 for missing items — that is expected
  const res = http.get('https://api.example.com/items/nonexistent');

  check(res, {
    'returns 404 for missing item': (r) => r.status === 404,
    'error message is descriptive': (r) => r.json('error') !== '',
  });
}
` + "```" + `

## Data Management and Parameterization

### Use SharedArray for Large Datasets

` + "`SharedArray`" + ` shares data across VUs, saving memory:

` + "```javascript" + `
import { SharedArray } from 'k6/data';
import http from 'k6/http';

// Data is loaded once and shared across all VUs (read-only)
const users = new SharedArray('users', function () {
  return JSON.parse(open('./data/users.json'));
});

export default function () {
  const user = users[Math.floor(Math.random() * users.length)];
  http.post('https://api.example.com/login', JSON.stringify({
    username: user.username,
    password: user.password,
  }), { headers: { 'Content-Type': 'application/json' } });
}
` + "```" + `

### Parameterize with Environment Variables

Use ` + "`__ENV`" + ` for runtime configuration:

` + "```javascript" + `
const BASE_URL = __ENV.BASE_URL || 'https://test.k6.io';
const API_KEY = __ENV.API_KEY;

export default function () {
  const res = http.get(BASE_URL + '/api/data', {
    headers: { 'X-API-Key': API_KEY },
  });
}
` + "```" + `

Run with:
` + "```bash" + `
k6 run -e BASE_URL=https://staging.example.com -e API_KEY=secret script.js
` + "```" + `

### Use Execution Context for Unique Data

Avoid data collisions between VUs using the execution API:

` + "```javascript" + `
import exec from 'k6/execution';
import { SharedArray } from 'k6/data';

const users = new SharedArray('users', function () {
  return JSON.parse(open('./data/users.json'));
});

export default function () {
  // Each VU gets a unique user based on its ID
  const user = users[exec.vu.idInTest % users.length];
}
` + "```" + `

## Authentication Strategies

### Authenticate Once in setup()

Do not authenticate on every iteration — authenticate once and pass the token:

` + "```javascript" + `
import http from 'k6/http';

export function setup() {
  const res = http.post('https://api.example.com/auth/token', JSON.stringify({
    client_id: __ENV.CLIENT_ID,
    client_secret: __ENV.CLIENT_SECRET,
  }), { headers: { 'Content-Type': 'application/json' } });

  return { token: res.json('access_token') };
}

export default function (data) {
  http.get('https://api.example.com/protected', {
    headers: { Authorization: 'Bearer ' + data.token },
  });
}
` + "```" + `

### Per-VU Authentication When Needed

If each VU needs its own session, authenticate in the init or default function:

` + "```javascript" + `
import http from 'k6/http';
import exec from 'k6/execution';
import { SharedArray } from 'k6/data';

const credentials = new SharedArray('creds', function () {
  return JSON.parse(open('./data/credentials.json'));
});

export default function () {
  const cred = credentials[exec.vu.idInTest % credentials.length];
  const loginRes = http.post('https://api.example.com/login', JSON.stringify(cred), {
    headers: { 'Content-Type': 'application/json' },
  });
  const token = loginRes.json('token');

  // Use the token for subsequent requests in this iteration
  http.get('https://api.example.com/me', {
    headers: { Authorization: 'Bearer ' + token },
  });
}
` + "```" + `

## Monitoring and Observability

### Use Custom Metrics for Business KPIs

Track domain-specific metrics alongside HTTP metrics:

` + "```javascript" + `
import { Trend, Counter, Rate } from 'k6/metrics';
import http from 'k6/http';
import { check } from 'k6';

const orderDuration = new Trend('order_processing_time');
const orderCount = new Counter('orders_placed');
const orderSuccess = new Rate('order_success_rate');

export default function () {
  const res = http.post('https://api.example.com/orders', JSON.stringify({
    items: [{ id: 1, quantity: 2 }],
  }), { headers: { 'Content-Type': 'application/json' } });

  orderDuration.add(res.timings.duration);
  const success = res.status === 201;
  orderSuccess.add(success);
  if (success) {
    orderCount.add(1);
  }
}
` + "```" + `

### Tag Requests for Granular Analysis

Use tags to break down metrics by endpoint or operation:

` + "```javascript" + `
import http from 'k6/http';

export default function () {
  http.get('https://api.example.com/users', {
    tags: { name: 'GetUsers', type: 'list' },
  });

  http.get('https://api.example.com/users/1', {
    tags: { name: 'GetUser', type: 'detail' },
  });
}
` + "```" + `

## Design Patterns: Ramping, Stages, and Scenarios

### Use Scenarios for Realistic Workloads

Scenarios let you model different traffic patterns running simultaneously:

` + "```javascript" + `
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  scenarios: {
    browse: {
      executor: 'constant-vus',
      vus: 50,
      duration: '5m',
      exec: 'browsePage',
    },
    checkout: {
      executor: 'ramping-arrival-rate',
      startRate: 1,
      timeUnit: '1s',
      preAllocatedVUs: 20,
      maxVUs: 100,
      stages: [
        { duration: '2m', target: 10 },
        { duration: '3m', target: 10 },
        { duration: '1m', target: 0 },
      ],
      exec: 'checkoutFlow',
    },
  },
};

export function browsePage() {
  http.get('https://ecommerce.example.com/');
  sleep(2);
}

export function checkoutFlow() {
  http.post('https://ecommerce.example.com/cart/add', JSON.stringify({ item: 1 }), {
    headers: { 'Content-Type': 'application/json' },
  });
  sleep(1);
  http.post('https://ecommerce.example.com/checkout');
}
` + "```" + `

### Choose the Right Executor

| Executor | Use case |
|----------|----------|
| ` + "`shared-iterations`" + ` | Fixed total iterations split across VUs — good for one-time batch jobs |
| ` + "`per-vu-iterations`" + ` | Each VU runs a fixed number of iterations — good for per-user workflows |
| ` + "`constant-vus`" + ` | Constant number of VUs — simplest steady-state test |
| ` + "`ramping-vus`" + ` | Ramp VUs up/down — classic load profile |
| ` + "`constant-arrival-rate`" + ` | Fixed request rate regardless of response time — good for SLO testing |
| ` + "`ramping-arrival-rate`" + ` | Ramp request rate up/down — find breaking points |

### Ramp Up and Down Gracefully

Avoid slamming your system with full load from the start:

` + "```javascript" + `
export const options = {
  stages: [
    { duration: '2m', target: 50 },   // ramp up
    { duration: '5m', target: 50 },   // steady state
    { duration: '2m', target: 100 },  // push higher
    { duration: '5m', target: 100 },  // steady at peak
    { duration: '3m', target: 0 },    // ramp down
  ],
};
` + "```" + `

## Code Quality: Modules and Shared Code

### Extract Reusable Code into Modules

Keep test scripts clean by extracting helpers:

` + "```javascript" + `
// helpers/api.js
import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'https://api.example.com';

export function apiGet(path, token) {
  const res = http.get(BASE_URL + path, {
    headers: { Authorization: 'Bearer ' + token },
  });
  return res;
}

export function apiPost(path, body, token) {
  const res = http.post(BASE_URL + path, JSON.stringify(body), {
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + token,
    },
  });
  return res;
}

export function checkStatus(res, expectedStatus, name) {
  check(res, {
    [name || 'status is ' + expectedStatus]: (r) => r.status === expectedStatus,
  });
}
` + "```" + `

` + "```javascript" + `
// test.js — uses the helper module
import { apiGet, apiPost, checkStatus } from './helpers/api.js';

export function setup() {
  const res = apiPost('/auth/login', { user: 'admin', pass: 'admin' });
  return { token: res.json('token') };
}

export default function (data) {
  const res = apiGet('/users', data.token);
  checkStatus(res, 200, 'get users');
}
` + "```" + `

### Use options Exports for Configuration

Keep configuration in the test script (or import from a shared config):

` + "```javascript" + `
// config/load-profile.js
export const smokeTest = {
  vus: 1,
  duration: '30s',
  thresholds: {
    http_req_duration: ['p(95)<500'],
  },
};

export const loadTest = {
  stages: [
    { duration: '5m', target: 100 },
    { duration: '10m', target: 100 },
    { duration: '5m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<800', 'p(99)<1500'],
    http_req_failed: ['rate<0.01'],
  },
};
` + "```" + `

## Browser Testing Best Practices

### Use k6 Browser for End-to-End Testing

Combine protocol-level and browser-level testing:

` + "```javascript" + `
import { browser } from 'k6/browser';
import { check } from 'k6';

export const options = {
  scenarios: {
    ui: {
      executor: 'shared-iterations',
      options: {
        browser: {
          type: 'chromium',
        },
      },
    },
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    await page.goto('https://test.k6.io/');

    const header = await page.locator('h1');
    check(await header.textContent(), {
      'header is correct': (text) => text.includes('Welcome'),
    });

    await page.locator('a[href="/contacts.php"]').click();
    await page.waitForNavigation();

    check(page, {
      'navigated to contacts': (p) => p.url().includes('/contacts'),
    });
  } finally {
    await page.close();
  }
}
` + "```" + `

### Mix Browser and Protocol Tests

Run browser tests alongside API tests for comprehensive coverage:

` + "```javascript" + `
import { browser } from 'k6/browser';
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    api_load: {
      executor: 'constant-vus',
      vus: 50,
      duration: '5m',
      exec: 'apiTest',
    },
    browser_flow: {
      executor: 'constant-vus',
      vus: 2,
      duration: '5m',
      exec: 'browserTest',
      options: {
        browser: {
          type: 'chromium',
        },
      },
    },
  },
};

export function apiTest() {
  const res = http.get('https://test.k6.io/api/data');
  check(res, { 'api status 200': (r) => r.status === 200 });
  sleep(1);
}

export async function browserTest() {
  const page = await browser.newPage();
  try {
    await page.goto('https://test.k6.io/');
    check(page, {
      'page loaded': (p) => p.url() === 'https://test.k6.io/',
    });
    sleep(3);
  } finally {
    await page.close();
  }
}
` + "```" + `

### Keep Browser VU Counts Low

Browser tests consume significantly more resources than protocol tests. Use 2-5 browser VUs for realistic frontend testing, and rely on protocol-level VUs for load generation:

` + "```javascript" + `
export const options = {
  scenarios: {
    // Heavy load via protocol
    protocol: {
      executor: 'ramping-vus',
      stages: [
        { duration: '5m', target: 200 },
        { duration: '10m', target: 200 },
        { duration: '5m', target: 0 },
      ],
      exec: 'protocolTest',
    },
    // Light browser presence for Web Vitals and UX metrics
    browser: {
      executor: 'constant-vus',
      vus: 3,
      duration: '20m',
      exec: 'browserTest',
      options: { browser: { type: 'chromium' } },
    },
  },
};
` + "```" + `

## Summary Checklist

- [ ] Use the k6 lifecycle stages correctly (init, setup, default, teardown)
- [ ] Set meaningful thresholds for CI/CD gating
- [ ] Use checks, not console.log, for validations
- [ ] Add realistic think time with sleep()
- [ ] Use SharedArray for large datasets
- [ ] Tag requests for granular metric analysis
- [ ] Extract reusable code into modules
- [ ] Choose the right executor for your use case
- [ ] Ramp up gradually — do not start at peak load
- [ ] Keep browser VUs low, use protocol VUs for load
- [ ] Use custom metrics for business-specific KPIs
- [ ] Authenticate once in setup() when possible
- [ ] Parameterize environment-specific values with __ENV
`

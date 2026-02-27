package docs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"go.k6.io/k6/cmd/state"
	"golang.org/x/term"
)

func newCmd(gs *state.GlobalState) *cobra.Command {
	return newDocsCmd(gs, os.Stdout)
}

func newDocsCmd(gs *state.GlobalState, stdout io.Writer) *cobra.Command {
	var opts docsOpts

	cmd := &cobra.Command{
		Use:   "docs [topic] [subtopic...]",
		Short: "Print k6 documentation",
		Long:  "Access k6 documentation from the command line.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDocs(gs, cmd, args, &opts)
		},
	}

	cmd.Flags().BoolVar(&opts.list, "list", false, "List subtopics instead of showing content")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Print all documentation")
	cmd.PersistentFlags().StringVar(&opts.version, "version", "", "Override k6 version for docs lookup")
	cmd.PersistentFlags().StringVar(&opts.cacheDir, "cache-dir", "", "Override cache directory")

	searchCmd := &cobra.Command{
		Use:   "search <term>",
		Short: "Search documentation",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args, &opts)
		},
	}
	cmd.AddCommand(searchCmd)

	return cmd
}

type docsOpts struct {
	list     bool
	all      bool
	version  string
	cacheDir string
}

func runSearch(cmd *cobra.Command, args []string, opts *docsOpts) error {
	version, cacheDir, idx, err := setup(opts.version, opts.cacheDir)
	if err != nil {
		return err
	}

	term := strings.Join(args, " ")
	printSearch(cmd.OutOrStdout(), idx, term, cacheDir, version)
	return nil
}

func runDocs(gs *state.GlobalState, cmd *cobra.Command, args []string, opts *docsOpts) error {
	version, cacheDir, idx, err := setup(opts.version, opts.cacheDir)
	if err != nil {
		return err
	}

	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	logMode(gs, isTTY)

	cfg, cfgErr := loadConfig()
	if cfgErr != nil && gs != nil {
		gs.Logger.Warnf("docs: ignoring invalid config: %v", cfgErr)
	}

	baseW := cmd.OutOrStdout()
	var buf *bytes.Buffer
	w := io.Writer(baseW)

	if cfg.Renderer != "" && isTTY {
		buf = &bytes.Buffer{}
		w = buf
	}

	if opts.all {
		printAll(w, idx, cacheDir, version)
		return pipeRenderer(buf, os.Stdout, baseW, cfg.Renderer)
	}

	if opts.list && len(args) == 0 {
		printTopLevelList(w, idx)
		return pipeRenderer(buf, os.Stdout, baseW, cfg.Renderer)
	}

	if len(args) == 0 {
		printTOC(w, idx, version)
		return pipeRenderer(buf, os.Stdout, baseW, cfg.Renderer)
	}

	if args[0] == "best-practices" {
		if err := printBestPractices(w, cacheDir, version); err != nil {
			return err
		}
		return pipeRenderer(buf, os.Stdout, baseW, cfg.Renderer)
	}

	slug := ResolveWithLookup(args, func(s string) bool {
		_, ok := idx.Lookup(s)
		return ok
	})

	sec, ok := idx.Lookup(slug)
	if !ok {
		return fmt.Errorf("topic not found: %s", strings.Join(args, " "))
	}

	if opts.list {
		printList(w, idx, slug)
		return pipeRenderer(buf, os.Stdout, baseW, cfg.Renderer)
	}

	printSection(w, idx, sec, cacheDir, version)
	return pipeRenderer(buf, os.Stdout, baseW, cfg.Renderer)
}

func logMode(gs *state.GlobalState, isTTY bool) {
	if gs == nil {
		return
	}
	if isTTY {
		gs.Logger.Debug("docs: interactive mode (stdout is TTY)")
	} else {
		gs.Logger.Debug("docs: agent mode (stdout is not a TTY)")
	}
}

func pipeRenderer(buf *bytes.Buffer, stdout, fallback io.Writer, renderer string) error {
	if buf == nil || buf.Len() == 0 {
		return nil
	}

	raw := buf.Bytes()

	parts := strings.Fields(renderer)
	if len(parts) == 0 {
		_, err := fallback.Write(raw)
		return err
	}

	rc := exec.Command(parts[0], parts[1:]...) //nolint:gosec // user-configured renderer
	rc.Stdin = bytes.NewReader(raw)
	rc.Stdout = stdout
	rc.Stderr = os.Stderr

	if err := rc.Run(); err != nil {
		_, writeErr := fallback.Write(raw)
		return writeErr
	}

	return nil
}

// setup resolves the version, ensures docs are cached, and loads the index.
// It checks flags, then env vars, then auto-detection for both version and
// cache directory.
func setup(versionFlag, cacheDirFlg string) (version, cacheDir string, idx *Index, err error) {
	version = versionFlag
	if version == "" {
		version = os.Getenv("K6_DOCS_VERSION")
	}
	if version == "" {
		version, err = DetectK6Version()
		if err != nil {
			return "", "", nil, fmt.Errorf("detect k6 version: %w", err)
		}
	}

	cacheDir = cacheDirFlg
	if cacheDir == "" {
		cacheDir = os.Getenv("K6_DOCS_CACHE_DIR")
	}

	if cacheDir == "" {
		cacheDir, err = EnsureDocs(version, http.DefaultClient)
		if err != nil {
			return "", "", nil, fmt.Errorf("ensure docs: %w", err)
		}
	}

	idx, err = LoadIndex(cacheDir)
	if err != nil {
		return "", "", nil, fmt.Errorf("load index: %w", err)
	}

	return version, cacheDir, idx, nil
}

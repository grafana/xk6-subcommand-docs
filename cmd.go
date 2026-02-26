package docs

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.k6.io/k6/cmd/state"
)

func newCmd(gs *state.GlobalState) *cobra.Command {
	var (
		listFlag    bool
		allFlag     bool
		versionFlag string
		cacheDirFlg string
	)

	cmd := &cobra.Command{
		Use:   "docs [topic] [subtopic...]",
		Short: "Print k6 documentation",
		Long:  "Access k6 documentation from the command line.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			version, cacheDir, idx, err := setup(versionFlag, cacheDirFlg)
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()

			if allFlag {
				printAll(w, idx, cacheDir, version)
				return nil
			}

			if len(args) == 0 {
				printTOC(w, idx, version)
				return nil
			}

			// Special case: "best-practices" as first arg.
			if args[0] == "best-practices" {
				return printBestPractices(w, cacheDir)
			}

			slug := Resolve(args)

			sec, ok := idx.Lookup(slug)
			if !ok {
				return fmt.Errorf("topic not found: %s", strings.Join(args, " "))
			}

			if listFlag {
				printList(w, idx, slug)
				return nil
			}

			printSection(w, idx, sec, cacheDir, version)
			return nil
		},
	}

	cmd.Flags().BoolVar(&listFlag, "list", false, "List subtopics instead of showing content")
	cmd.Flags().BoolVar(&allFlag, "all", false, "Print all documentation")
	cmd.PersistentFlags().StringVar(&versionFlag, "version", "", "Override k6 version for docs lookup")
	cmd.PersistentFlags().StringVar(&cacheDirFlg, "cache-dir", "", "Override cache directory")

	searchCmd := &cobra.Command{
		Use:   "search <term>",
		Short: "Search documentation",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, cacheDir, idx, err := setup(versionFlag, cacheDirFlg)
			if err != nil {
				return err
			}

			term := strings.Join(args, " ")
			printSearch(cmd.OutOrStdout(), idx, term, cacheDir)
			return nil
		},
	}
	cmd.AddCommand(searchCmd)

	return cmd
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

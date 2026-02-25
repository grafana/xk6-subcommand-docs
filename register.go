// Package docs contains the xk6-subcommand-docs extension.
package docs

import (
	"github.com/spf13/cobra"
	"go.k6.io/k6/cmd/state"
	"go.k6.io/k6/subcommand"
)

func init() {
	subcommand.RegisterExtension("docs", newCmd)
}

func newCmd(gs *state.GlobalState) *cobra.Command {
	// Stub for now — will be implemented later.
	return &cobra.Command{
		Use:   "docs",
		Short: "Print k6 documentation",
	}
}

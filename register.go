// Package docs contains the xk6-subcommand-docs extension.
package docs

import "go.k6.io/k6/subcommand"

func init() {
	subcommand.RegisterExtension("docs", newCmd)
}

package search

import (
	"github.com/growthbook/gb-find-code-refs/internal/helpers"
	"github.com/growthbook/gb-find-code-refs/options"
)

// Get a list of delimiters to use for flag key matching
// If defaults are disabled, only additional configured delimiters will be used
func GetDelimiters(opts options.Options) []string {
	delims := []string{`"`, `'`, "`"}
	if opts.Delimiters.DisableDefaults {
		delims = []string{}
	}

	delims = append(delims, opts.Delimiters.Additional...)

	return helpers.Dedupe(delims)
}

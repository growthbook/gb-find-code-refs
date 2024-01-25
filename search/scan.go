package search

import (
	"github.com/growthbook/gb-find-code-refs/flags"
	"github.com/growthbook/gb-find-code-refs/internal/gb"
	"github.com/growthbook/gb-find-code-refs/internal/log"
	"github.com/growthbook/gb-find-code-refs/options"
)

// Scan checks the configured directory for flags based on the options configured for Code References.
func Scan(opts options.Options, dir string) (Matcher, []gb.ReferenceHunksRep) {
	flagKeys := flags.GetFlagKeys(opts)
	matcher := NewMultiProjectMatcher(opts, dir, flagKeys)

	refs, err := SearchForRefs(dir, matcher)
	if err != nil {
		log.Error.Fatalf("error searching for flag key references: %s", err)
	}

	return matcher, refs
}

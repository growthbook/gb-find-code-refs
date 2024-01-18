package search

import (
	"github.com/launchdarkly/ld-find-code-refs/v2/flags"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/gb"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/log"
	"github.com/launchdarkly/ld-find-code-refs/v2/options"
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

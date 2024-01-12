package flags

import (
	"fmt"
	"os"

	"github.com/launchdarkly/ld-find-code-refs/v2/internal/helpers"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/ld"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/log"
	"github.com/launchdarkly/ld-find-code-refs/v2/options"
)

const (
	minFlagKeyLen = 3 // Minimum flag key length helps reduce the number of false positives
)

func GetFlagKeys(opts options.Options, repoParams ld.RepoParams) []string {
	ignoreServiceErrors := opts.IgnoreServiceErrors
	flags, err := getFlags()
	if err != nil {
		helpers.FatalServiceError(fmt.Errorf("could not parse flag keys: %w", err), ignoreServiceErrors)
	}

	filteredFlags, omittedFlags := filterShortFlagKeys(flags)
	if len(filteredFlags) == 0 {
		log.Info.Printf("no flag keys longer than the minimum flag key length (%v) were found, exiting early",
			minFlagKeyLen)
		os.Exit(0)
	} else if len(omittedFlags) > 0 {
		log.Warning.Printf("omitting %d flags with keys less than minimum (%d) for project: %s", len(omittedFlags), minFlagKeyLen)
	}
	return filteredFlags
}

// Very short flag keys lead to many false positives when searching in code,
// so we filter them out.
func filterShortFlagKeys(flags []string) (filtered []string, omitted []string) {
	filteredFlags := []string{}
	omittedFlags := []string{}
	for _, flag := range flags {
		if len(flag) >= minFlagKeyLen {
			filteredFlags = append(filteredFlags, flag)
		} else {
			omittedFlags = append(omittedFlags, flag)
		}
	}
	return filteredFlags, omittedFlags
}

func getFlags() ([]string, error) {
	// TODO read from stdin or file (path provided as arg)
	return []string{"feature-list-realtime-graphs", "github-integration"}, nil
}

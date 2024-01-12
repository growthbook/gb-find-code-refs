package coderefs

import (
	"strings"

	"github.com/launchdarkly/ld-find-code-refs/v2/internal/git"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/helpers"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/ld"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/log"
	"github.com/launchdarkly/ld-find-code-refs/v2/internal/validation"
	"github.com/launchdarkly/ld-find-code-refs/v2/options"
	"github.com/launchdarkly/ld-find-code-refs/v2/search"
)

func Run(opts options.Options, output bool) {
	absPath, err := validation.NormalizeAndValidatePath(opts.Dir)
	if err != nil {
		log.Error.Fatalf("could not validate directory option: %s", err)
	}

	log.Info.Printf("absolute directory path: %s", absPath)

	branchName := opts.Branch
	revision := opts.Revision
	var gitClient *git.Client
	var commitTime int64
	if revision == "" {
		gitClient, err = git.NewClient(absPath, branchName, opts.AllowTags)
		if err != nil {
			log.Error.Fatalf("%s", err)
		}
		branchName = gitClient.GitBranch
		revision = gitClient.GitSha
		commitTime = gitClient.GitTimestamp
	}

	repoParams := ld.RepoParams{
		Type:              opts.RepoType,
		Name:              opts.RepoName,
		Url:               opts.RepoUrl,
		CommitUrlTemplate: opts.CommitUrlTemplate,
		HunkUrlTemplate:   opts.HunkUrlTemplate,
		DefaultBranch:     opts.DefaultBranch,
	}

	matcher, refs := search.Scan(opts, repoParams, absPath)

	var updateId *int
	if opts.UpdateSequenceId >= 0 {
		updateIdOption := opts.UpdateSequenceId
		updateId = &updateIdOption
	}

	branch := ld.BranchRep{
		Name:             strings.TrimPrefix(branchName, "refs/heads/"),
		Head:             revision,
		UpdateSequenceId: updateId,
		SyncTime:         helpers.MakeTimestamp(),
		References:       refs,
		CommitTime:       commitTime,
	}

	if output {
		generateHunkOutput(opts, matcher, branch, repoParams)
	}

	if gitClient != nil {
		runExtinctions(opts, matcher, branch, repoParams, gitClient)
	}
}

func deleteStaleBranches(ldApi ld.ApiClient, repoName string, remoteBranches map[string]bool) error {
	branches, err := ldApi.GetCodeReferenceRepositoryBranches(repoName)
	if err != nil {
		return err
	}

	staleBranches := calculateStaleBranches(branches, remoteBranches)
	if len(staleBranches) > 0 {
		log.Debug.Printf("marking stale branches for code reference pruning: %v", staleBranches)
		err = ldApi.PostDeleteBranchesTask(repoName, staleBranches)
		if err != nil {
			return err
		}
	}

	return nil
}

func calculateStaleBranches(branches []ld.BranchRep, remoteBranches map[string]bool) []string {
	staleBranches := []string{}
	for _, branch := range branches {
		if !remoteBranches[branch.Name] {
			staleBranches = append(staleBranches, branch.Name)
		}
	}
	log.Info.Printf("found %d stale branches to be marked for code reference pruning", len(staleBranches))
	return staleBranches
}

func generateHunkOutput(opts options.Options, matcher search.Matcher, branch ld.BranchRep, repoParams ld.RepoParams) {
	outDir := opts.OutDir

	if outDir != "" {
		outPath, err := branch.WriteToCSV(outDir, "default", repoParams.Name, opts.Revision)
		if err != nil {
			log.Error.Fatalf("error writing code references to csv: %s", err)
		}
		log.Info.Printf("wrote code references to %s", outPath)
	}

	if opts.Debug {
		branch.PrintReferenceCountTable()
	}

	totalFlags := 0
	for _, searchElems := range matcher.Elements {
		totalFlags += len(searchElems.Elements)
	}
	log.Info.Printf(
		"found %d code references across %d flags and %d files",
		branch.TotalHunkCount(),
		totalFlags,
		len(branch.References),
	)
}

func runExtinctions(opts options.Options, matcher search.Matcher, branch ld.BranchRep, repoParams ld.RepoParams, gitClient *git.Client) {
	if opts.Lookback > 0 {
		var removedFlags []ld.ExtinctionRep

		flagCounts := branch.CountByProjectAndFlag(matcher.GetElements(), []string{"default"})
		missingFlags := []string{}
		for flag, count := range flagCounts["default"] {
			if count == 0 {
				missingFlags = append(missingFlags, flag)
			}
		}
		log.Info.Printf("checking if %d flags without references were removed in the last %d commits for project: %s", len(missingFlags), opts.Lookback, "default")
		removedFlagsByProject, err := gitClient.FindExtinctions(missingFlags, matcher, opts.Lookback+1)
		if err != nil {
			log.Warning.Printf("unable to generate flag extinctions: %s", err)
		} else {
			log.Info.Printf("found %d removed flags", len(removedFlagsByProject))
		}
		removedFlags = append(removedFlags, removedFlagsByProject...)

		// TODO replace this with way to output to stdout instead
		// if len(removedFlags) > 0 && !opts.DryRun {
		//   err := ldApi.PostExtinctionEvents(removedFlags, repoParams.Name, branch.Name)
		//   if err != nil {
		//     log.Error.Printf("error sending extinction events to LaunchDarkly: %s", err)
		//   }
		// }
	}
}

package coderefs

import (
	"strings"

	"github.com/growthbook/gb-find-code-refs/internal/gb"
	"github.com/growthbook/gb-find-code-refs/internal/git"
	"github.com/growthbook/gb-find-code-refs/internal/helpers"
	"github.com/growthbook/gb-find-code-refs/internal/log"
	"github.com/growthbook/gb-find-code-refs/internal/validation"
	"github.com/growthbook/gb-find-code-refs/options"
	"github.com/growthbook/gb-find-code-refs/search"
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

	matcher, refs := search.Scan(opts, absPath)

	var updateId *int
	if opts.UpdateSequenceId >= 0 {
		updateIdOption := opts.UpdateSequenceId
		updateId = &updateIdOption
	}

	branch := gb.BranchRep{
		Name:             strings.TrimPrefix(branchName, "refs/heads/"),
		Head:             revision,
		UpdateSequenceId: updateId,
		SyncTime:         helpers.MakeTimestamp(),
		References:       refs,
		CommitTime:       commitTime,
	}

	if output {
		generateHunkOutput(opts, matcher, branch)
	}

	if gitClient != nil {
		runExtinctions(opts, matcher, branch, gitClient)
	}
}

func generateHunkOutput(opts options.Options, matcher search.Matcher, branch gb.BranchRep) {
	// default to current directory
	outDir := opts.OutDir
	if outDir == "" {
		outDir = "."
	}

	outPath, err := branch.WriteToJSON(outDir, "default", opts.Revision)
	if err != nil {
		log.Error.Fatalf("error writing code references to csv: %s", err)
	}
	log.Info.Printf("wrote code references to %s", outPath)

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

func runExtinctions(opts options.Options, matcher search.Matcher, branch gb.BranchRep, gitClient *git.Client) {
	if opts.Lookback > 0 {
		var removedFlags []gb.ExtinctionRep

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

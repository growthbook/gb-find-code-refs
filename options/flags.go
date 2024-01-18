package options

type flag struct {
	name         string
	short        string
	defaultValue interface{}
	usage        string
}

// Options that are available as command line flags
// TODO update to reflect changes
var flags = []flag{
	{
		name:         "allowTags",
		defaultValue: false,
		usage:        "Enables storing references for tags. The tag will be listed as a branch.",
	},
	{
		name:         "branch",
		short:        "b",
		defaultValue: "",
		usage: `The currently checked out branch. If not provided, branch
name will be auto-detected. Provide this option when using CI systems that
leave the repository in a detached HEAD state.`,
	},
	{
		name:         "commitUrlTemplate",
		defaultValue: "",
		usage: `If provided, LaunchDarkly will attempt to generate links to
your VCS service provider per commit.
Example: https://github.com/launchdarkly/gb-find-code-refs/commit/${sha}.
Allowed template variables: 'branchName', 'sha'. If "commitUrlTemplate" is not provided, but "repoUrl" is provided and "repoType" is not custom, LaunchDarkly will attempt to automatically generate source code links for the given "repoType".`,
	},
	{
		name:         "contextLines",
		short:        "C",
		defaultValue: 2,
		usage: `The number of context lines to send to LaunchDarkly. If < 0, no
source code will be sent to LaunchDarkly. If 0, only the lines containing
flag references will be sent. If > 0, will send that number of context
lines above and below the flag reference. A maximum of 5 context lines
may be provided.`,
	},
	{
		name:         "debug",
		defaultValue: false,
		usage:        "Enables verbose debug logging",
	},
	{
		name:         "defaultBranch",
		short:        "B",
		defaultValue: "main",
		usage: `The default branch. The LaunchDarkly UI will default to this branch.
If not provided, will fallback to 'main'.`,
	},
	{
		name:         "dir",
		short:        "d",
		defaultValue: "",
		usage:        "Path to existing checkout of the repository.",
	},
	{
		name:         "hunkUrlTemplate",
		defaultValue: "",
		usage: `If provided, LaunchDarkly will attempt to generate links to 
your VCS service provider per code reference. 
Example: https://github.com/launchdarkly/gb-find-code-refs/blob/${sha}/${filePath}#L${lineNumber}.
Allowed template variables: 'sha', 'filePath', 'lineNumber'. If "hunkUrlTemplate" is not provided, but "repoUrl" is provided and "repoType" is not custom, LaunchDarkly will attempt to automatically generate source code links for the given "repoType".`,
	},
	{
		name:         "ignoreServiceErrors",
		short:        "i",
		defaultValue: false,
		usage: `If enabled, the scanner will terminate with exit code 0 when the
LaunchDarkly API is unreachable or returns an unexpected response.`,
	},
	{
		name:         "lookback",
		short:        "l",
		defaultValue: 10,
		usage: `Sets the number of git commits to search in history for
whether a feature flag was removed from code. May be set to 0 to disabled this feature. Setting this option to a high value will increase search time.`,
	},
	{
		name:         "outDir",
		short:        "o",
		defaultValue: "",
		usage: `If provided, will output a csv file containing all code references for
the project to this directory.`,
	},
	{
		name:         "repoType",
		short:        "T",
		defaultValue: "custom",
		usage: `The repo service provider. Used to correctly categorize repositories in the
LaunchDarkly UI. Acceptable values: bitbucket|custom|github|gitlab.`,
	},
	{
		name:         "repoUrl",
		short:        "u",
		defaultValue: "",
		usage:        `The URL for the repository. If provided and "repoType" is not custom, LaunchDarkly will attempt to automatically generate source code links for the given "repoType".`,
	},
	{
		name:         "revision",
		short:        "R",
		defaultValue: "",
		usage:        `Use this option to scan non-git codebases. The current revision of the repository to be scanned. If set, the version string for the scanned repository will not be inferred, and branch garbage collection will be disabled. The "branch" option is required when "revision" is set.`,
	},
	{
		name:         "updateSequenceId",
		short:        "s",
		defaultValue: -1,
		usage: `An integer representing the order number of code reference updates.
Used to version updates across concurrent executions of the flag finder.
If not provided, data will always be updated. If provided, data will
only be updated if the existing "updateSequenceId" is less than the new
"updateSequenceId". Examples: the time a "git push" was initiated, CI
build number, the current unix timestamp.`,
	},
	{
		name:         "userAgent",
		defaultValue: "",
		usage:        `(Internal) Platform where code references is run.`,
	},
	{
		name:         "flagsPath",
		short:        "f",
		defaultValue: "",
		usage:        "Required path to a JSON file containing a list of flag keys (array of strings). The scanner will search for references to the flags in this file.",
	},
}

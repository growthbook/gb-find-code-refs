package options

type flag struct {
	name         string
	short        string
	defaultValue interface{}
	usage        string
}

// Options that are available as command line flags
var flags = []flag{
	{
		name:         "allowTags",
		defaultValue: false,
		usage:        "Enables parsing references for tags.",
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
		name:         "contextLines",
		short:        "C",
		defaultValue: 2,
		usage:        `The number of context lines to include with each code reference. If 0, only the lines containing flag references will be sent. If > 0, will include that number of context lines above and below the flag reference. A maximum of 5 context lines may be provided. (default 2)`,
	},
	{
		name:         "debug",
		defaultValue: false,
		usage:        "Enables verbose debug logging",
	},
	{
		name:         "dir",
		short:        "d",
		defaultValue: "",
		usage:        "Path to existing checkout of the repository.",
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
		usage:        `If provided, will output the JSON file containing all code references to this directory. Otherwise, will output JSON file to current working directory.`,
	},
	{
		name:         "revision",
		short:        "R",
		defaultValue: "",
		usage:        `Use this option to scan non-git codebases. The current revision of the repository to be scanned. If set, the version string for the scanned repository will not be inferred, and branch garbage collection will be disabled. The "branch" option is required when "revision" is set.`,
	},
	{
		name:         "flagsPath",
		short:        "f",
		defaultValue: "",
		usage:        "Required path to a JSON file containing a list of flag keys (array of strings). The scanner will search for references to the flags in this file.",
	},
}

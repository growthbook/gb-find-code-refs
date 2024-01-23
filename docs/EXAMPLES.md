# Examples

The section provides examples of various `bash` commands to execute `gb-find-code-refs` (when installed in the system PATH) with various configurations. We recommend reading through the following examples to gain an understanding of common configurations. For more information on advanced configuration, see [CONFIGURATION.md](CONFIGURATION.md)

## Basic configuration

```bash
gb-find-code-refs \
  --dir="/path/to/git/repo" \
  --flagsPath="/path/to/flags.json"
```

## Configuration with context lines

https://docs.launchdarkly.com/home/code/code-references#configuring-context-lines

```bash
gb-find-code-refs \
  --dir="/path/to/git/repo" \
  --flagsPath="/path/to/flags.json" \
  --contextLines=3 # can be up to 5. If < 0, no source code will be sent to LD
```

## Scanning non-git repositories

By default, `gb-find-code-refs` will attempt to infer repository metadata from a git configuration. If you are scanning a codebase with a version control system other than git, you must use the `--revision` and `--branch` options to manually provide information about your codebase.

```bash
gb-find-code-refs \
  --dir="/path/to/git/repo" \
  --flagsPath="/path/to/flags.json" \
  --revision="REPO_REVISION_STRING" \ # e.g. a version hash
  --branch="dev"
```

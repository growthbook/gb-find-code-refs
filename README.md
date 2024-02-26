# gb-find-code-refs

Command line program for generating flag code references.

## Execution via CLI

The command line program may be run manually, and executed in an environment of your choosing. The program requires your `git` repo to be cloned locally, and the currently checked out branch will be scanned for code references.

We recommend incorporating `gb-find-code-refs` into your CI/CD build process. `gb-find-code-refs` should run whenever a commit is pushed to your repository.

## Example usage

```bash
# run CLI utility against codebase with feature flags provided in flags.json, repo name set to growthbook/growthbook
$ ./gb-find-code-refs -d ../growthbook -f ../flags.json -n growthbook/growthbook
INFO: 2024/02/26 10:42:42 coderefs.go:25: absolute directory path: /Users/alice/dev/growthbook/growthbook
INFO: 2024/02/26 10:42:42 git.go:47: git branch: main
INFO: 2024/02/26 10:42:42 coderefs.go:71: wrote code references to /Users/alice/dev/gb-find-code-refs/coderefs_main.json
INFO: 2024/02/26 10:42:42 coderefs.go:81: found 12 code references across 7 flags and 6 files

# post results to an endpoint, such as growthbook's code references endpoint
$ curl -XPOST -H "Authorization: Bearer ..." -H "Content-Type: application/json" your-growthbook-host/api/v1/code-refs -d @coderefs_main.json
# successful response
{
  "featuresUpdated": [
    "onboarding-banner",
    "new-checkout-flow",
    "extra-red-button"
  ]
}%
```

### Prerequisites

If you are scanning a git repository, `gb-find-code-refs` requires git (tested with version 2.21.0) to be installed on the system path.

### Installing

#### Docker

`gb-find-code-refs` is available as a [docker image](https://hub.docker.com/repository/docker/growthbook/gb-find-code-refs/general). The image provides an entrypoint for `gb-find-code-refs`, to which command line arguments may be passed. If using the entrypoint, your repository to be scanned should be mounted as a volume. Otherwise, you may override the entrypoint and access `gb-find-code-refs` directly from the shell.

```bash
docker pull growthbook/gb-find-code-refs
docker run \
  -v /path/to/your/repo:/repo \
  growthbook/gb-find-code-refs \
  --dir="/repo"
  --flagsPath="./flags.json"
```

### Configuration

`gb-find-code-refs` provides a number of configuration options to customize how code references are generated.

-   [All configuration options are documented in CONFIGURATION.md](docs/CONFIGURATION.md)
-   [Common configuration examples are documented in EXAMPLES.md](docs/EXAMPLES.md)
-   [Detailed information on configuring feature flag aliases is documented in ALIASES.md](docs/ALIASES.md)

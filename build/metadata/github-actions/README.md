# GrowthBook Code References with GitHub Actions

This GitHub Action is a utility that automatically populates code references in GrowthBook. This is useful for finding references to feature flags in your code, both for reference and for code cleanup.

## Configuration

Create a new Actions workflow in your selected GitHub repository (e.g. `code-references.yml`) in the `.github/workflows` directory of your repository. Under "Edit new file", paste the following code:

```yaml
on: push
name: Find feature flag code references
concurrency:
    group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
    cancel-in-progress: true

jobs:
    growthBookCodeReferences:
        name: GrowthBook Code References
        runs-on: ubuntu-latest
        steps:
            - uses: actions/checkout@v4
              with:
                  fetch-depth: 11 # This value must be set if the lookback configuration option is defined for find-code-refs. Read more: https://github.com/growthbook/gb-find-code-refs#searching-for-unused-flags-extinctions
            - name: GrowthBook Code References
              uses: growthbook/gb-find-code-refs@v2.11.5
              with:
                  flagsPath: GB_FLAGS_PATH
```

We strongly recommend that you update the second `uses` attribute value to reference the latest tag in the [growthbook/gb-find-code-refs repository](https://github.com/growthbook/gb-find-code-refs). This will pin your workflow to a particular version of the `growthbook/gb-find-code-refs` action. Also, make sure to change `projKey` to the key of the GrowthBook project associated with this repository.

Commit this file under a new branch. Submit as a PR to your code reviewers to be merged into your default branch. You do not need to have this branch merged into the default branch for code references to appear in the GrowthBook UI for your flags; code references will appear for this newly created branch.

As shown in the above example, the workflow should run on the `push` event, and contain an action provided by the [growthbook/gb-find-code-refs repository](https://github.com/growthbook/gb-find-code-refs). The environment variables should be included as a secret.

## Additional Examples

The example below is the same as first, but it also excludes any `dependabot` branches. We suggest excluding any automatically generated branches where flags do not change.

```yaml
on:
    push:
        branches-ignore:
            - "dependabot/**"

name: Find GrowthBook flag code references
concurrency:
    group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
    cancel-in-progress: true

jobs:
    growthBookCodeReferences:
        name: GrowthBook Code References
        runs-on: ubuntu-latest
        steps:
            - uses: actions/checkout@v4
              with:
                  fetch-depth: 11 # This value must be set if the lookback configuration option is not disabled for find-code-references. Read more: https://github.com/growthbook/gb-find-code-refs#searching-for-unused-flags-extinctions
            - name: GrowthBook Code References
              uses: growthbook/gb-find-code-refs@v2.11.5
              with:
                  flagsPath: GB_FLAGS_PATH
```

## Troubleshooting

Once your workflow has been created, the best way to confirm that the workflow is executing correctly is to create a new pull request with the workflow file and verify that the newly created action succeeds.

If the action fails, there may be a problem with your configuration. To investigate, dig into the action's logs to view any error messages.

<!-- action-docs-inputs -->

## Inputs

| parameter    | description                                                                                                                                                                                                                                                                                                                                                                                                                     | required | default |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| allowTags    | Enable storing references for tags. Lists the tag as a branch.                                                                                                                                                                                                                                                                                                                                                                  | `false`  | false   |
| contextLines | The number of context lines above and below a code reference for the job to send to GrowthBook. By default, the flag finder will not send any context lines to GrowthBook. If < 0, it will send no source code to GrowthBook. If 0, it will send only the lines containing flag references. If > 0, it will send that number of context lines above and below the flag reference. You may provide a maximum of 5 context lines. | `false`  | 2       |
| debug        | Enable verbose debug logging.                                                                                                                                                                                                                                                                                                                                                                                                   | `false`  | false   |
| lookback     | Set the number of commits to search in history for whether you removed a feature flag from code. You may set to 0 to disable this feature. Setting this option to a high value will increase search time.                                                                                                                                                                                                                       | `false`  | 10      |

<!-- action-docs-inputs -->

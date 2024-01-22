# Note: These commands pertain to the development of gb-find-code-refs.
#       They are not intended for use by the end-users of this program.
SHELL=/bin/bash
GORELEASER_VERSION=v1.20.0

build:
	go build ./cmd/...

init:
	pre-commit install

test: lint
	go test ./...

lint:
	pre-commit run -a --verbose golangci-lint

# Generate docs about GitHub Action inputs and updates README.md
github-action-docs:
	cd build/metadata/github-actions && npx action-docs -u --no-banner

# Strip debug informatino from production builds
BUILD_FLAGS = -ldflags="-s -w"

compile-macos-binary:
	GOOS=darwin GOARCH=amd64 go build ${BUILD_FLAGS} -o out/gb-find-code-refs ./cmd/gb-find-code-refs

compile-windows-binary:
	GOOS=windows GOARCH=amd64 go build ${BUILD_FLAGS} -o out/gb-find-code-refs.exe ./cmd/gb-find-code-refs

compile-linux-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/cmd/gb-find-code-refs ./cmd/gb-find-code-refs

compile-github-actions-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/github-actions/gb-find-code-refs-github-action ./build/package/github-actions

# Get the lines added to the most recent changelog update (minus the first 2 lines)
RELEASE_NOTES=<(GIT_EXTERNAL_DIFF='bash -c "diff --unchanged-line-format=\"\" $$2 $$5" || true' git log --ext-diff -1 --pretty= -p CHANGELOG.md)

echo-release-notes:
	@cat $(RELEASE_NOTES)

define publish_docker
	test $(1) || (echo "Please provide tag"; exit 1)
	docker build -t launchdarkly/$(3):$(1) build/package/$(4)
	docker push launchdarkly/$(3):$(1)
	# test $(2) && (echo "Not pushing latest tag for prerelease")
	test $(2) || docker tag launchdarkly/$(3):$(1) launchdarkly/$(3):latest
	test $(2) || docker push launchdarkly/$(3):latest
endef

clean:
	rm -rf out/
	rm -f build/pacakge/cmd/gb-find-code-refs
	rm -f build/package/github-actions/gb-find-code-refs-github-action

RELEASE_CMD=curl -sL https://git.io/goreleaser | GOPATH=$(mktemp -d) VERSION=$(GORELEASER_VERSION) GITHUB_TOKEN=$(GITHUB_TOKEN) bash -s -- --clean --release-notes $(RELEASE_NOTES)

publish:
	$(RELEASE_CMD)

test-publish:
	curl -sL https://git.io/goreleaser | VERSION=$(GORELEASER_VERSION) bash -s -- --clean --skip-publish --skip-validate

products-for-release:
	$(RELEASE_CMD) --skip-publish --skip-validate

.PHONY: init test lint compile-github-actions-binary compile-macos-binary compile-linux-binary compile-windows-binary echo-release-notes publish-all clean build

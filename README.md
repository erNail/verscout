# verscout

<div align="center">
  <img src="./docs/img/verscout-icon.png" width="150" alt="verscout icon">
  <p>
    Find the latest version tag, calculate the next version, print to STDOUT
    - no tagging, no bumping, no changelog, no publishing.
  </p>

  [![Go Report Card](https://goreportcard.com/badge/github.com/erNail/verscout)](https://goreportcard.com/report/github.com/erNail/verscout)
  [![Release](https://img.shields.io/github/v/release/erNail/verscout)](https://github.com/erNail/verscout/releases/latest)
  [![License](https://img.shields.io/github/license/erNail/verscout)](LICENSE)
</div>

## What is `verscout`?

`verscout` is a single binary CLI tool. `verscout latest` will print the latest version tag of your repository to STDOUT.
`verscout next` will calculate the next version tag based on
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/), and print it to STDOUT.
If no latest or next version exists, `verscout` will print nothing to STDOUT and exit with code `0` by default.

`verscout` will not create and push any tags.
It will not bump any versions in any files.
It will not create a changelog.
It will not publish any artifacts.

## Why `verscout`?

I wanted to have a lightweight tool that can tell me
the latest and next version based on Git tags and
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for use in my CI processes.
Most other tools known to me did not fulfill my requirements:

- Print the latest version to STDOUT
- Print the next version to STDOUT
- Print no next version to STDOUT if there are no new commits that should cause a version bump
  and exit with code `0`

## Getting Started

### Install `verscout`

#### Via Homebrew

```shell
brew install erNail/tap/verscout
```

#### Via Binary

Check the [releases](https://github.com/erNail/verscout/releases) for the available binaries.
Download the correct binary and add it to your `$PATH`.

#### Via Go

```shell
go install github.com/erNail/verscout
```

#### Via Container

```shell
docker pull ernail/verscout:<LATEST_GITHUB_RELEASE_VERSION>
```

#### From Source

Check out this repository and run the following:

```shell
go build
```

Add the resulting binary to your `$PATH`.

### Run `verscout`

#### Get the latest version tag

```shell
verscout latest
```

For verscout to find the latest version, the tags need to be in the format `vMAJOR.MINOR.PATCH`
or `MAJOR.MINOR.PATCH`

#### Calculate the next version

```shell
verscout next
```

For verscout to calculate the next version,
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
using `fix`, `feat` or `BREAKING CHANGE:` need to exist since the latest version tag.

### Configure `verscout`

To get a complete list of the configuration options, please use the `--help` or `-h` flag.

```shell
verscout --help
```

#### Global options

##### Working Directory

By default, `verscout` will run in the current working directory.
Use the `--dir` flag to change this behavior.

```shell
verscout --dir ./my-other-repository
```

#### Options for `verscout latest`

##### Exit Code if no latest version is found

By default, `verscout latest` will exit with code `0` if no latest version is found due to expected reasons.

The expected reasons are:

- There are no existing tags
- There are no valid version tags

You can change this behavior with the `--exit-code` flag

```shell
verscout latest --exit-code 4
```

#### Options for `verscout next`

By default, `verscout next` will exit with code `0` if no next version is found due to expected reasons.

The expected reasons are:

- There are no new commit messages since the last tag
- There are no new commit message since the last tag that use any keywords that will cause a version bump

You can change this behavior with the `--exit-code` flag

```shell
verscout next --exit-code 4
```

### Limitations

- The format of the version tags is currently not configurable
- Which [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
  should cause a version bump is currently not configurable
- Bumping any prerelease or release candidate versions is not supported

## Planned Features and Limitations

Please check the open [GitHub Issues](https://github.com/erNail/homebrew-tap/issues)
to get an overview of the planned features.

## Development

### Dependencies

Please check the tasks in the [`taskfile.yaml`](./taskfile.yaml) for any tools you might need.

### Testing

```shell
task test
```

### Linting

```shell
task lint
```

### Running

```shell
task run -- --help
```

### Building

```shell
task build
```

### Building Container Images

```shell
task build-image
```

### Test GitHub Actions

```shell
task test-github-actions
```

### Test Release

```shell
task release-test
```

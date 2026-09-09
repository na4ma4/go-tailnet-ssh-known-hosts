# go-tailnet-ssh-known-hosts

SSH host key program for verifying Tailscale tailnet hosts.

## Installation

```sh
go install github.com/na4ma4/go-tailnet-ssh-known-hosts/cmd/_ssh_tailnet_knownhosts@latest
```

Homebrew users can install the command from the formula in this repository:

```sh
brew tap na4ma4/tap
brew ssh_tailnet_knownhosts
```

Published release binaries are available on the [Releases](https://github.com/na4ma4/go-tailnet-ssh-known-hosts/releases) page.

The command requires an active local Tailscale connection. It connects directly
to the Tailscale LocalAPI and does not invoke the Tailscale CLI.

## Usage

```sh
_ssh_tailnet_knownhosts HOST [...HOST]
```

The command prints the SSH host keys advertised by matching online peers in
OpenSSH `known_hosts` format. Hosts can be DNS names or Tailscale IP addresses.
Matching DNS names is case-insensitive and ignores a trailing dot. Results are
sorted and duplicate entries are removed. If no online peer matches, the
command prints nothing.

## OpenSSH configuration

Add it to `~/.ssh/config`:

```sshconfig
Host *
    KnownHostsCommand /absolute/path/to/_ssh_tailnet_knownhosts %H %k
```

`%H` passes the destination hostname to the command, `%k` sends the host alias.
Make sure to use its absolute path in `KnownHostsCommand`.

## Development

```sh
mage
mage build
mage runE "examplehostname"
```

`mage` runs formatting, dependency updates, linting, and the race-enabled test
suite.

`mage build` writes a local binary to `artifacts/bin/`. `mage release` runs
GoReleaser locally; releases created in GitHub Actions use the same
`.goreleaser.yml` configuration.

## Releases

Commits merged into `main` are collected by Release Please. It opens a release
pull request with the changelog and version update; merging that pull request
creates the Git tag and GitHub release. The release workflow then uses
GoReleaser to build archives for Linux, macOS, and Windows on amd64 and arm64.

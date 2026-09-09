# go-tailnet-ssh-known-hosts

SSH host key program for verifying Tailscale tailnet hosts.

## Installation

```sh
go install github.com/na4ma4/go-tailnet-ssh-known-hosts/cmd/_ssh_tailnet_knownhosts@latest
```

The command requires the Tailscale CLI and an active local Tailscale connection.
It reads `tailscale status --json` and returns host keys advertised by online
peers in OpenSSH `known_hosts` format.

Add it to `~/.ssh/config`:

```sshconfig
Host *
    KnownHostsCommand _ssh_tailnet_knownhosts %H %k
```

OpenSSH calls the command with the destination host and requested key type.
The command prints no output when that host/key type is not advertised by an
online tailnet peer.

## Development

```sh
go test ./...
go build ./cmd/_ssh_tailnet_knownhosts
```

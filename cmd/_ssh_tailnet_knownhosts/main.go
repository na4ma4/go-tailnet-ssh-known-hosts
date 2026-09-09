// _ssh_tailnet_knownhosts prints Tailscale SSH host keys for OpenSSH.
package main

import (
	"context"
	"fmt"
	"os"

	"tailscale.com/client/local"

	"github.com/na4ma4/go-tailnet-ssh-known-hosts/internal/knownhosts"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: _ssh_tailnet_knownhosts HOST [...HOST]")
		os.Exit(2)
	}

	status, err := (&local.Client{}).Status(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entries, err := knownhosts.Keys(status, os.Args[1:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

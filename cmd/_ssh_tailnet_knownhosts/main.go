// _ssh_tailnet_knownhosts prints Tailscale SSH host keys for OpenSSH.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/na4ma4/go-tailnet-ssh-known-hosts/internal/knownhosts"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: _ssh_tailnet_knownhosts HOST KEYTYPE")
		os.Exit(2)
	}

	status, err := knownhosts.CommandOutput(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entries, err := knownhosts.Keys(status, os.Args[1], os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

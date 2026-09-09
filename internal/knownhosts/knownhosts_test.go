package knownhosts_test

import (
	"net/netip"
	"testing"

	"github.com/na4ma4/go-tailnet-ssh-known-hosts/internal/knownhosts"

	"tailscale.com/ipn/ipnstate"
	"tailscale.com/types/key"
)

func TestKeys(t *testing.T) {
	t.Parallel()

	status := &ipnstate.Status{
		Peer: map[key.NodePublic]*ipnstate.PeerStatus{
			{}: {
				DNSName:      "node.example.ts.net.",
				TailscaleIPs: []netip.Addr{netip.MustParseAddr("100.100.100.10")},
				Online:       true,
				SSH_HostKeys: []string{"ssh-ed25519 AAAA-ed25519", "ssh-rsa AAAA-rsa"},
			},
		},
	}

	entries, err := knownhosts.Keys(status, "NODE.EXAMPLE.TS.NET")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"NODE.EXAMPLE.TS.NET ssh-ed25519 AAAA-ed25519",
		"NODE.EXAMPLE.TS.NET ssh-rsa AAAA-rsa",
	}
	if len(entries) != len(want) || entries[0] != want[0] || entries[1] != want[1] {
		t.Fatalf("Keys() = %#v, want %#v", entries, want)
	}
}

func TestKeysMatchesIP(t *testing.T) {
	t.Parallel()

	status := &ipnstate.Status{
		Peer: map[key.NodePublic]*ipnstate.PeerStatus{
			{}: {
				TailscaleIPs: []netip.Addr{netip.MustParseAddr("100.100.100.10")},
				Online:       true,
				SSH_HostKeys: []string{"ssh-ed25519 AAAA"},
			},
		},
	}

	entries, err := knownhosts.Keys(status, "100.100.100.10")
	if err != nil || len(entries) != 1 {
		t.Fatalf("Keys() = %#v, %v; want one entry", entries, err)
	}
}

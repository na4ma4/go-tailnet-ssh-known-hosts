// Package knownhosts finds SSH host keys advertised by Tailscale peers.
package knownhosts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"os/exec"
	"sort"
	"strings"
)

// Status is the part of tailscale status --json that this package needs.
type Status struct {
	Peers map[string]PeerStatus `json:"Peer"`
}

// PeerStatus contains the peer's addresses and advertised SSH host keys.
type PeerStatus struct {
	DNSName      string   `json:"DNSName"`
	TailscaleIPs []string `json:"TailscaleIPs"`
	Online       bool     `json:"Online"`
	SSHHostKeys  []string `json:"SSHHostKeys"`
}

// CommandOutput runs tailscale and returns its JSON status.
func CommandOutput(ctx context.Context) ([]byte, error) {
	output, err := exec.CommandContext(ctx, "tailscale", "status", "--json").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("tailscale status --json: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return output, nil
}

// Keys returns known_hosts entries for host and keyType from online peers.
func Keys(statusJSON []byte, hosts ...string) ([]string, error) {
	// if keyType == "" {
	// 	return nil, errors.New("key type is required")
	// }
	if len(hosts) == 0 {
		return nil, errors.New("at least one host is required")
	}

	var status Status
	if err := json.Unmarshal(statusJSON, &status); err != nil {
		return nil, fmt.Errorf("decode tailscale status: %w", err)
	}

	entries := make([]string, 0)
	seen := make(map[string]struct{})
	for _, peer := range status.Peers {
		if !peer.Online {
			continue
		}
		for _, host := range hosts {
			if !matchesHost(peer, host) {
				continue
			}
			for _, hostKey := range peer.SSHHostKeys {
				fields := strings.Fields(hostKey)
				if len(fields) < 2 {
					continue
				}
				entry := host + " " + strings.Join(fields, " ")
				if _, ok := seen[entry]; ok {
					continue
				}
				seen[entry] = struct{}{}
				entries = append(entries, entry)
			}
		}
	}

	sort.Strings(entries)

	return entries, nil
}

func matchesHost(peer PeerStatus, host string) bool {
	if equalDNSName(peer.DNSName, host) {
		return true
	}

	var requestedIP netip.Addr
	{
		var err error
		requestedIP, err = netip.ParseAddr(host)
		if err != nil {
			return false
		}
	}

	for _, address := range peer.TailscaleIPs {
		if ip, err := netip.ParseAddr(address); err == nil && ip == requestedIP {
			return true
		}
	}

	return false
}

func equalDNSName(a, b string) bool {
	if strings.EqualFold(strings.TrimSuffix(a, "."), strings.TrimSuffix(b, ".")) {
		return true
	}

	if strings.HasPrefix(a, b) {
		return true
	}

	return false
}

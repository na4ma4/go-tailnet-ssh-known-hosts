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
func Keys(statusJSON []byte, host, keyType string) ([]string, error) {
	if host == "" || keyType == "" {
		return nil, errors.New("host and key type are required")
	}

	var status Status
	if err := json.Unmarshal(statusJSON, &status); err != nil {
		return nil, fmt.Errorf("decode tailscale status: %w", err)
	}

	entries := make([]string, 0)
	seen := make(map[string]struct{})
	for _, peer := range status.Peers {
		if !peer.Online || !matchesHost(peer, host) {
			continue
		}
		for _, hostKey := range peer.SSHHostKeys {
			fields := strings.Fields(hostKey)
			if len(fields) < 2 || fields[0] != keyType {
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

	sort.Strings(entries)
	return entries, nil
}

func matchesHost(peer PeerStatus, host string) bool {
	if equalDNSName(peer.DNSName, host) {
		return true
	}
	requestedIP, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	for _, address := range peer.TailscaleIPs {
		if ip, err := netip.ParseAddr(address); err == nil && ip == requestedIP {
			return true
		}
	}
	return false
}

func equalDNSName(a, b string) bool {
	return strings.EqualFold(strings.TrimSuffix(a, "."), strings.TrimSuffix(b, "."))
}

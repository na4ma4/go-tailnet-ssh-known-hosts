// Package knownhosts finds SSH host keys advertised by Tailscale peers.
package knownhosts

import (
	"errors"
	"net/netip"
	"slices"
	"sort"
	"strings"

	"tailscale.com/ipn/ipnstate"
)

// Keys returns known_hosts entries for host and keyType from online peers.
func Keys(status *ipnstate.Status, hosts ...string) ([]string, error) {
	// if keyType == "" {
	// 	return nil, errors.New("key type is required")
	// }
	if status == nil {
		return nil, errors.New("status is required")
	}
	if len(hosts) == 0 {
		return nil, errors.New("at least one host is required")
	}

	entries := make([]string, 0)
	seen := make(map[string]struct{})
	for _, peer := range status.Peer {
		for _, entry := range peerEntries(peer, hosts) {
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

func peerEntries(peer *ipnstate.PeerStatus, hosts []string) []string {
	if peer == nil || !peer.Online {
		return nil
	}

	entries := make([]string, 0)
	for _, host := range hosts {
		if !matchesHost(peer, host) {
			continue
		}
		for _, hostKey := range peer.SSH_HostKeys {
			fields := strings.Fields(hostKey)
			if len(fields) < 2 {
				continue
			}
			entries = append(entries, host+" "+strings.Join(fields, " "))
		}
	}

	return entries
}

func matchesHost(peer *ipnstate.PeerStatus, host string) bool {
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

	return slices.Contains(peer.TailscaleIPs, requestedIP)
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

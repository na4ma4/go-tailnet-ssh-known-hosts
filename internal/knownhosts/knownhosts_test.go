package knownhosts

import "testing"

func TestKeys(t *testing.T) {
	status := []byte(`{
		"Peer": {
			"node": {
				"DNSName": "node.example.ts.net.",
				"TailscaleIPs": ["100.100.100.10"],
				"Online": true,
				"SSHHostKeys": ["ssh-ed25519 AAAA-ed25519", "ssh-rsa AAAA-rsa"]
			},
			"offline": {
				"DNSName": "offline.example.ts.net.",
				"Online": false,
				"SSHHostKeys": ["ssh-ed25519 AAAA-offline"]
			}
		}
	}`)

	entries, err := Keys(status, "NODE.EXAMPLE.TS.NET", "ssh-ed25519")
	if err != nil {
		t.Fatal(err)
	}
	if want := "NODE.EXAMPLE.TS.NET ssh-ed25519 AAAA-ed25519"; len(entries) != 1 || entries[0] != want {
		t.Fatalf("Keys() = %#v, want %#v", entries, []string{want})
	}
}

func TestKeysMatchesIP(t *testing.T) {
	status := []byte(`{"Peer":{"node":{"TailscaleIPs":["100.100.100.10"],"Online":true,"SSHHostKeys":["ssh-ed25519 AAAA"]}}}`)
	entries, err := Keys(status, "100.100.100.10", "ssh-ed25519")
	if err != nil || len(entries) != 1 {
		t.Fatalf("Keys() = %#v, %v; want one entry", entries, err)
	}
}

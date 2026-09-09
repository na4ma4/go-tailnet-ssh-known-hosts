class SshTailnetKnownhosts < Formula
  desc "SSH known-hosts command for Tailscale tailnet peers"
  homepage "https://github.com/na4ma4/go-tailnet-ssh-known-hosts"
  url "https://github.com/na4ma4/go-tailnet-ssh-known-hosts/archive/refs/heads/main.tar.gz"
  version "0.1.0"
  head "https://github.com/na4ma4/go-tailnet-ssh-known-hosts.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/_ssh_tailnet_knownhosts"
  end

  test do
    assert_match "usage:", shell_output("#{bin}/ssh-tailnet-knownhosts 2>&1", 2)
  end
end

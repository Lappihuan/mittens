# Mittens

A kubectl plugin for intercepting HTTP traffic to Kubernetes Services using mitmproxy.

[![Build status][shield-build-status]][build-status]
[![Latest release][shield-latest-release]][latest-release]
[![License][shield-license]][license]

## Usage

```sh
kubectl mittens SERVICE [OPTIONS]
```

**Examples:**
```sh
kubectl mittens my-service -n my-namespace          # Auto-detect port
kubectl mittens my-service -p 8080                  # Explicit port
kubectl mittens my-service -p 443 --https           # HTTPS service
```

**Options:**
- `-n, --namespace STRING`: Target namespace
- `-p, --port INT`: Service port (auto-detected if omitted)
- `--https`: Enable for HTTPS services
- `-i, --image STRING`: Custom proxy image
- `--command-args STRING`: Custom mitmproxy arguments

**What happens:**
1. Deploy mitmproxy sidecar to service pods
2. Redirect traffic through mitmproxy
3. Open interactive mitmproxy TUI
4. Auto-cleanup on exit (Ctrl+C)

## Installation

**Binary:** Download from [Releases](https://github.com/Lappihuan/mittens/releases)

**From source:** `go install github.com/Lappihuan/mittens/cmd/kubectl-mittens@latest`

**With Krew:** `kubectl krew install mittens`

## K9s Integration

Add to `~/.k9s/plugins.yaml`:

```yaml
plugins:
  mittens:
    shortCut: Ctrl-T
    description: "mittens: inject mitmproxy sidecar"
    scopes:
      - services
    command: kubectl
    background: false
    args:
      - mittens
      - $NAME
      - -n
      - $NAMESPACE
```

## License

Apache 2.0. See [LICENSE](LICENSE) and [ATTRIBUTION.md](ATTRIBUTION.md).

## Attribution

Mittens is a fork of [kubetap](https://github.com/soluble-ai/kubetap) by Soluble Inc, originally licensed under Apache 2.0. Mittens maintains full Apache 2.0 compliance and acknowledges the original project. See [ATTRIBUTION.md](ATTRIBUTION.md) for full attribution details.

[shield-build-status]: https://github.com/Lappihuan/mittens/workflows/build/badge.svg?branch=master
[shield-latest-release]: https://img.shields.io/github/v/release/Lappihuan/mittens?include_prereleases&label=release&sort=semver
[shield-license]: https://img.shields.io/github/license/Lappihuan/mittens.svg
[license]: https://github.com/Lappihuan/mittens/blob/master/LICENSE
[latest-release]: https://github.com/Lappihuan/mittens/releases
[build-status]: https://github.com/Lappihuan/mittens/actions

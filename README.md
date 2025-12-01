# Mittens

Mittens is a kubectl plugin that enables operators to easily deploy interactive mitmproxy instances for Kubernetes Services. Named after the German word "Mitten" (in the middle), Mittens places mitmproxy directly in front of Kubernetes Services for real-time HTTP traffic interception and inspection.

[![Build status][shield-build-status]][build-status]
[![Latest release][shield-latest-release]][latest-release]
[![Go version][shield-go-version]][go-version]
[![License][shield-license]][license]

## What is Mittens?

Mittens is a kubectl plugin that streamlines the process of deploying a mitmproxy sidecar to intercept and inspect HTTP traffic destined for Kubernetes Services. Unlike traditional approaches that require manual deployment manifests or complex proxy configurations, Mittens automates the entire process with simple commands.

## Why Mittens?

### The Problem

When debugging microservices in Kubernetes, operators often need to:

1. Inspect HTTP traffic to understand service interactions
2. Test error handling and edge cases
3. Debug production issues without code changes
4. Examine inter-service communication

Traditional solutions require:

- Manual sidecar injection into deployment manifests
- Service routing reconfiguration
- Complex port forwarding setups
- Browser-based proxy interfaces

### The Solution

Mittens automates sidecar deployment and provides direct, interactive access to mitmproxy's Terminal User Interface (TUI). Users simply enable mittens for a service, interact with mitmproxy in the terminal, and automatic cleanup happens when they exit.

## Features

- **One-Command Deployment**: Deploy mitmproxy with a single kubectl command
- **Interactive TUI Access**: Direct terminal access to mitmproxy for real-time inspection
- **Automatic Cleanup**: Sidecars and configurations are automatically removed when you exit
- **K9s Integration**: Quick mittens toggle directly from the k9s dashboard
- **Context-Aware**: Respects kubectl context, namespace, and service selection
- **No Port Forwarding**: Direct `kubectl exec` for clean, simple architecture
- **Zero Configuration**: Works out of the box with sensible defaults

## Installation

### Binary Release

Download binaries for macOS, Linux, and Windows from the [Releases page](https://github.com/Lappihuan/mittens/releases).

### From Source

```sh
go install github.com/Lappihuan/mittens/cmd/kubectl-mittens@latest
```

Or clone and build locally:

```sh
git clone https://github.com/Lappihuan/mittens.git
cd mittens
go install ./cmd/kubectl-mittens
```

### With Krew

```sh
kubectl krew install mittens
```

## Quick Start

### Enable mittens for a Service

```sh
kubectl mittens on my-service -n my-namespace -p 8080
```

This will:
1. Deploy a mitmproxy sidecar to pods matching the service
2. Redirect traffic through mitmproxy
3. Open an interactive tmux session with mitmproxy TUI
4. Automatically clean up when you exit

### Interact with mitmproxy

Once connected:

```
- Arrow keys: Navigate through requests
- q: Quit mitmproxy
- Tab: Cycle through request/response views
- ?: View help
- Ctrl+B then D: Detach from tmux (keep mitmproxy running)
```

### Disable mittens

Simply exit the mitmproxy session, and mittens automatically cleans up:

```
- Ctrl+C or 'q' in mitmproxy will exit and trigger cleanup
- All sidecar containers are removed
- Service routing is restored
- ConfigMaps are deleted
```

## Usage

### Tap On - Enable mittens

```sh
kubectl mittens on SERVICE [OPTIONS]
```

**Options:**
- `-n, --namespace STRING`: Target namespace (default: current context)
- `-p, --port INT`: Target service port
- `--https`: Enable for HTTPS services
- `-i, --image STRING`: Custom proxy image
- `--command-args STRING`: Custom mitmproxy arguments

**Examples:**

```sh
# Tap HTTP service on port 8080
kubectl mittens on my-api -n default -p 8080

# Tap HTTPS service
kubectl mittens on my-api -n default -p 443 --https

# Use custom image
kubectl mittens on my-api -p 8080 -i ghcr.io/custom/mitmproxy:latest

# Custom mitmproxy mode
kubectl mittens on my-api -p 8080 --command-args "mitmweb"
```

### Tap Off - Disable mittens

```sh
kubectl mittens off SERVICE [OPTIONS]
```

This is rarely needed as cleanup happens automatically, but useful for emergency cleanup:

```sh
kubectl mittens off my-service -n my-namespace
```

### List - Show Active Mittens

```sh
kubectl mittens list [OPTIONS]
```

**Options:**
- `-n, --namespace STRING`: Filter by namespace (default: all namespaces)

**Example:**

```sh
# List all services with mittens enabled
kubectl mittens list

# List in specific namespace
kubectl mittens list -n production
```

## K9s Integration

Mittens integrates with [k9s](https://k9scli.io/) for quick access from the terminal UI dashboard.

### Installation

1. Copy the mittens plugin to your k9s config:

```sh
cat docs/k9s-plugin.yaml >> ~/.k9s/plugins.yaml
```

2. Restart k9s

### Usage

Navigate to **Services** or **Pods** view and press `Ctrl+M` to enable mittens for the selected service or pod.

For details, see [K9s Integration Guide](docs/getting_started/k9s-integration.md).

## Architecture

### Traditional Approach (Before)

```
User Browser
    |
    v (Port Forward 4000)
Localhost:4000 <--> mitmweb UI
    |
    v (HTTP Traffic)
Kubernetes Service
```

Issues:
- Requires port forwarding
- Needs browser access
- Manual cleanup
- Extra moving parts

### Mittens Approach (Now)

```
kubectl mittens on --> Sidecar Deployed
                   --> Direct kubectl exec
                   --> Interactive mitmproxy TUI
User Terminal <----> mitmproxy tmux session
                   --> User exits
                   --> Auto cleanup
```

Benefits:
- Direct terminal access
- No port forwarding
- Automatic cleanup
- Single process flow
- Minimal overhead

## How It Works

1. **Validation**: Mittens verifies the target service and namespace exist
2. **Injection**: A mitmproxy container is added to pods matching the service selector
3. **Configuration**: A ConfigMap is created with mitmproxy settings
4. **Service Patching**: The Service is patched to route traffic through mittens
5. **Attachment**: `kubectl exec` attaches to an interactive tmux session
6. **User Interaction**: User inspects traffic in mitmproxy TUI
7. **Exit**: When user exits, automatic cleanup removes all mittens components

## Requirements

- kubectl 1.19+
- Kubernetes 1.19+
- mitmproxy 12.0+ (runs in container)

## Differences from Original Kubetap

Mittens is a fork of the original [kubetap](https://github.com/soluble-ai/kubetap) project by Soluble, with significant architectural changes:

| Feature | Kubetap | Mittens |
|---------|---------|---------|
| **Proxy UI** | mitmweb (browser) | mitmproxy TUI (terminal) |
| **Access** | Port forwarding | Direct kubectl exec |
| **User Experience** | Browser interface | Terminal interface |
| **Cleanup** | Manual | Automatic |
| **CLI** | `kubectl tap` | `kubectl mittens` |
| **Architecture** | Complex | Simple |

See [ATTRIBUTION.md](ATTRIBUTION.md) for full details about the fork and original project attribution.

## Examples

### Debugging a Microservice

```sh
# Enable mittens for payment-service
kubectl mittens on payment-service -n production -p 8080

# In mitmproxy, watch for failed requests
# Inspect headers and request bodies
# Test error responses

# Exit when done - automatic cleanup happens
```

### Testing Rate Limiting

```sh
# Enable mittens for rate-limiter-service
kubectl mittens on rate-limiter -n staging -p 9000

# Use mitmproxy to see how requests are throttled
# Modify requests to test edge cases
# Inspect response headers for rate limit info
```

### Troubleshooting Payment Issues

```sh
# Monitor payment processor calls
kubectl mittens on payment-processor -n prod -p 443 --https

# Inspect all HTTPS requests in real-time
# Review exact request/response payloads
# Identify formatting or authentication issues
```

## Development

### Building from Source

```sh
git clone https://github.com/Lappihuan/mittens
cd mittens
go build -o mittens ./cmd/kubectl-mittens
```

### Running Tests

```sh
make test
```

### Installing Locally

```sh
go install ./cmd/kubectl-mittens
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

Mittens is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Attribution

Mittens is a fork of [kubetap](https://github.com/soluble-ai/kubetap) by Soluble Inc, originally licensed under Apache 2.0. Mittens maintains full Apache 2.0 compliance and acknowledges the original project. See [ATTRIBUTION.md](ATTRIBUTION.md) for full attribution details.

## Support

- Documentation: [https://github.com/Lappihuan/mittens](https://github.com/Lappihuan/mittens)
- Issues: [GitHub Issues](https://github.com/Lappihuan/mittens/issues)
- Original Project: [kubetap](https://github.com/soluble-ai/kubetap)

[shield-go-version]: https://img.shields.io/github/go-mod/go-version/Lappihuan/mittens
[shield-build-status]: https://github.com/Lappihuan/mittens/workflows/mittens/badge.svg?branch=master
[shield-latest-release]: https://img.shields.io/github/v/release/Lappihuan/mittens?include_prereleases&label=release&sort=semver
[shield-license]: https://img.shields.io/github/license/Lappihuan/mittens.svg
[license]: https://github.com/Lappihuan/mittens/blob/master/LICENSE
[go-version]: https://github.com/Lappihuan/mittens/blob/master/go.mod
[latest-release]: https://github.com/Lappihuan/mittens/releases
[build-status]: https://github.com/Lappihuan/mittens/actions
[kubectl-plugin]: https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/

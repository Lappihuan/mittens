# Installation

Mittens is a kubectl plugin that can be installed in several ways.

## With Krew (Recommended)

Mittens can be installed with [krew](https://github.com/kubernetes-sigs/krew), the kubectl plugin manager:

```sh
kubectl krew install mittens
```

Once installed, you can use mittens with:

```sh
kubectl mittens on <service>
kubectl mittens off <service>
kubectl mittens list
```

## From Binary Release

Pre-compiled binary releases for Mac (arm64/amd64), Linux (amd64/arm64), and Windows are available from the [GitHub Releases page](https://github.com/Lappihuan/mittens/releases).

1. Download the binary for your platform from the latest release
2. Extract the binary and move it to a directory in your `PATH` (e.g., `/usr/local/bin`)
3. Make the binary executable: `chmod +x mittens`
4. Move to `$HOME/.krew/bin/` to use as a kubectl plugin: `mv mittens ~/.krew/bin/kubectl-mittens`

Then use:

```sh
kubectl mittens on <service>
```

Or if not in `~/.krew/bin/`:

```sh
./mittens on <service>
```

## From Source

To build mittens from source, clone the repository and install:

```sh
git clone https://github.com/Lappihuan/mittens.git
cd mittens
go install ./cmd/kubectl-mittens
```

This will build the binary as `kubectl-mittens` and install it to `$GOPATH/bin`. To use it as a kubectl plugin, either:

1. Add `$GOPATH/bin` to your `PATH`, or
2. Create a symlink: `ln -s $GOPATH/bin/kubectl-mittens ~/.krew/bin/kubectl-mittens`

Then verify the installation:

```sh
kubectl mittens --help
```

## Using the Docker Image

Mittens can also run directly in your Kubernetes cluster via a Docker container. The pre-built images are available at `ghcr.io/lappihuan/mittens`.

To run mittens as a pod (useful in CI/CD or serverless environments):

```sh
kubectl run mittens-tool --image=ghcr.io/lappihuan/mittens:latest -- mittens on <service>
```

However, for most use cases, installing as a kubectl plugin using krew is recommended, as it provides seamless integration with your kubectl workflow.

## Requirements

- Kubernetes 1.19 or later
- kubectl 1.19 or later
- mitmproxy will be deployed as a sidecar container in your target pod

## Next Steps

- See [Quick Start](quick-start.md) for a 5-minute walkthrough
- See [Usage Guide](usage.md) for full command reference
- See [K9s Integration](k9s-integration.md) for using mittens from k9s


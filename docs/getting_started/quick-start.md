# Quick Start

Get Mittens up and running in 5 minutes.

## Installation

### Quick Install

```sh
# Using kubectl krew
kubectl krew install mittens

# Or download binary from releases
# https://github.com/Lappihuan/mittens/releases
```

### Verify Installation

```sh
kubectl mittens version
# Output: version: vX.Y.Z, commit: ..., built at: ...
```

## Your First Tap

### Step 1: Deploy a Sample Service

For this example, we'll use a simple service. If you already have a service, skip to Step 2.

```sh
# Create a simple nginx deployment
kubectl create deployment web --image=nginx
kubectl expose deployment web --port=80 --target-port=80
```

### Step 2: Enable Mittens

```sh
kubectl mittens on web -p 80
```

This will:
1. Deploy a mitmproxy sidecar container
2. Patch the service to route traffic through mittens
3. Open an interactive mitmproxy terminal session

### Step 3: Inspect Traffic

Once mittens starts, you'll see the mitmproxy TUI with requests flowing through.

**Keyboard shortcuts in mitmproxy:**
- Arrow Up/Down: Navigate through requests
- Enter: View request/response details
- q: Quit mitmproxy
- Tab: Switch between panels
- ?: Show all commands

### Step 4: Exit and Automatic Cleanup

```sh
# In mitmproxy, press 'q' to quit
# or Ctrl+C to exit the terminal

# Mittens automatically:
# - Removes the sidecar container
# - Restores service routing
# - Cleans up all configurations
```

Verify cleanup:

```sh
kubectl mittens list
# Should show no tapped services
```

## Common Tasks

### Tap an HTTPS Service

```sh
kubectl mittens on secure-api -n default -p 443 --https
```

### Tap a Service on a Custom Port

```sh
kubectl mittens on my-service -n production -p 9000
```

### View All Active Taps

```sh
kubectl mittens list
```

### Disable Mittens Early

```sh
kubectl mittens off web -n default
```

## Troubleshooting

### Service Not Found

**Solution:** Verify the service name and namespace:
```sh
kubectl get svc -n default
```

### Pod Not Ready

This is normal. Mittens waits up to 90 seconds. Check pod status:
```sh
kubectl get pods -n default
```

## Next Steps

- Read the [Usage Guide](usage.md) for detailed command reference
- Set up [K9s Integration](k9s-integration.md) for quick tapping from k9s dashboard

Happy debugging!

Once attached, you can:
- View live HTTP/HTTPS traffic flowing through the proxy
- Inspect request/response headers and bodies
- Modify requests before they reach the target service
- Navigate using arrow keys, press `?` for help

To detach from tmux without stopping the proxy, press `Ctrl+B` then `D`.

## Listing active taps

All active taps can be listed using the following command, which can be constrained
to a specific namespace with `-n`:

```sh
$ kubectl tap list
Tapped Namespace/Service:

argocd/argocd-server
```

## Untapping the service

Once we are finished, we can remove the proxy and revert our tap by
turning it off:

```sh
$ kubectl tap off -n argocd argocd-server
Untapped Service "argocd-server"

```

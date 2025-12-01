# Usage

Mittens binary is `kubectl-mittens`, making it invokable as `kubectl mittens`.

Mittens inherits kubectl options including: `--context`, `--namespace`, `--user`, `--as`, etc.

## Commands

### mittens on - Enable mittens for a Service

Enable mittens (deploy mitmproxy) for the target Service.

```sh
kubectl mittens on SERVICE [OPTIONS]
```

**Options:**
- `-n, --namespace STRING`: Target namespace (default: current context)
- `-p, --port INT`: Target service port (required)
- `--https`: Enable for HTTPS services
- `-i, --image STRING`: Custom mitmproxy image (default: ghcr.io/lappihuan/mittens-mitmproxy:latest)
- `--command-args STRING`: Custom mitmproxy arguments (default: "mitmproxy")
- `--protocol STRING`: Protocol type (default: "http")

**Examples:**

```sh
# Tap HTTP service on port 8080
kubectl mittens on my-api -p 8080

# Tap HTTPS service
kubectl mittens on secure-api -p 443 --https

# Tap in specific namespace
kubectl mittens on my-service -n production -p 8080

# Use custom image
kubectl mittens on my-service -p 8080 -i ghcr.io/custom/mitmproxy:v13

# Custom mitmproxy mode
kubectl mittens on my-service -p 8080 --command-args "mitmproxy --mode reverse"
```

**What happens:**
1. Mittens deploys a mitmproxy sidecar container to matching pods
2. Service traffic is redirected through mittens
3. Interactive mitmproxy TUI opens in your terminal
4. When you exit, automatic cleanup removes all mittens components

### mittens off - Disable mittens for a Service

Remove mittens and restore original service routing.

```sh
kubectl mittens off SERVICE [OPTIONS]
```

**Options:**
- `-n, --namespace STRING`: Target namespace (default: current context)

**Examples:**

```sh
# Remove mittens from my-service
kubectl mittens off my-service

# Remove from specific namespace
kubectl mittens off my-service -n production
```

**Note:** This is rarely needed as cleanup happens automatically when you exit the mitmproxy session.

### mittens list - Show All Active Mittens

List all services currently running mittens.

```sh
kubectl mittens list [OPTIONS]
```

**Options:**
- `-n, --namespace STRING`: Filter by specific namespace (default: all namespaces)

**Examples:**

```sh
# List mittens in all namespaces
kubectl mittens list

# List mittens in production namespace
kubectl mittens list -n production

# List with full output
kubectl mittens list -n default
```

**Output:**
```
Tapped Namespace/Service:

default/my-api
production/payment-service
staging/test-service
```

## Kubernetes Inheritance

Mittens respects standard kubectl options:

```sh
# Use specific context
kubectl mittens on my-service -p 8080 --context=production-cluster

# Use different namespace
kubectl mittens on my-service -p 8080 -n staging

# Use service account
kubectl mittens on my-service -p 8080 --as=admin
```

## Running in Container

Run mittens as a Pod in Kubernetes:

```sh
# Create mittens pod
docker run -v "${HOME}/.kube/:/root/.kube/:ro" \
  ghcr.io/lappihuan/mittens:latest \
  mittens on -n myns -p 8080 myservice
```

Mittens automatically detects ServiceAccount tokens mounted to containers.

## Interacting with mitmproxy

Once connected to the mitmproxy TUI:

### Navigation

- `Arrow Up/Down`: Move through request list
- `Arrow Left/Right`: Scroll horizontally  
- `Enter`: Expand selected request
- `Tab`: Switch between different views (requests, details, etc.)

### Inspection

- `Space`: Preview full request/response
- `e`: Edit request before sending
- `m`: Change request method
- `d`: Delete request
- `U`: Switch to untagged view

### Utilities

- `a`: Add tag to request
- `t`: Tag by name
- `?`: Show help and all commands
- `q`: Quit mitmproxy
- `Ctrl+C`: Force quit

### Tmux Integration

Once in mitmproxy:

- `Ctrl+B then D`: Detach from tmux (keeps mitmproxy running)
- `Ctrl+B then [`: Enter copy mode
- `Ctrl+B then ]`: Paste copied text

Reattach later:
```sh
kubectl exec -it <pod-name> -n <namespace> -c mittens -- tmux attach-session -t mitmproxy
```

## Common Workflows

### Debugging Failed Requests

```sh
# Enable mittens
kubectl mittens on my-api -n prod -p 8080

# In mitmproxy:
# 1. Look for requests with 5xx status codes
# 2. Select request and press 'e' to inspect
# 3. View response details for error messages
# 4. Press 'q' to exit and cleanup
```

### Testing Rate Limiting

```sh
# Enable mittens on rate-limited service
kubectl mittens on rate-limiter -p 9000

# In mitmproxy:
# 1. Send multiple rapid requests
# 2. Watch for 429 Too Many Requests responses
# 3. Inspect rate-limit headers
# 4. Exit to cleanup
```

### Validating Request Format

```sh
# Enable mittens
kubectl mittens on payment-api -p 443 --https

# In mitmproxy:
# 1. Browse to requests from other services
# 2. Review request headers and body
# 3. Verify Content-Type and format
# 4. Check authentication headers
# 5. Exit when done
```

## Troubleshooting

### Service not found

```
Error: the provided Service does not exist
```

Solution: Check the service exists:
```sh
kubectl get svc -n <namespace>
```

### Port doesn't exist

```
Error: the target Service does not have the provided port
```

Solution: Verify the service port:
```sh
kubectl get svc <service> -n <namespace> -o yaml | grep -A 5 "ports:"
```

### Pod not ready

```
Waiting for Pod containers to become ready...
```

This is normal - mittens waits up to 90 seconds. Check pod status:
```sh
kubectl get pods -n <namespace> | grep <service>
```

### Cleanup failed

To manually cleanup a failed tap:
```sh
kubectl mittens off <service> -n <namespace>
```

## Cleanup

Mittens automatically cleans up when you exit the session. Manual cleanup removes:
- Mitmproxy sidecar container
- ConfigMap with mitmproxy configuration
- Service port redirections
- Deployment patches

## Running in Container

It is possible to schedule mittens as a Pod in Kubernetes:

```sh
docker run -v "${HOME}/.kube/:/root/.kube/:ro" \
  ghcr.io/lappihuan/mittens:latest \
  mittens on -n namespace -p 8080 myservice
```

When run in cluster, mittens automatically uses mounted ServiceAccount tokens.
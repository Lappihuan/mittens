# K9s Plugin Integration - Mittens

Mittens integrates with [k9s](https://k9scli.io/) to provide quick access to tapping services directly from the k9s terminal UI.

Named "mittens" as a nod to k9s (dog 🐕) and mitmproxy, plus "Mitten" being the German word for "in the middle".

## Installation

1. Copy the plugin configuration to your k9s plugins file:

```bash
# If ~/.k9s/plugins.yaml exists, append the mittens plugins:
cat docs/k9s-plugin.yaml >> ~/.k9s/plugins.yaml

# Or, if you don't have plugins.yaml yet:
mkdir -p ~/.k9s
cp docs/k9s-plugin.yaml ~/.k9s/plugins.yaml
```

2. Restart k9s to load the plugins

## Usage

### In Services View

Navigate to the **Services** view in k9s (usually with `:svc` or `:services`):

| Action | Key | Effect |
|--------|-----|--------|
| **Tap with Mittens** | `Ctrl+M` | Deploy mitmproxy sidecar to intercept traffic |

### In Pods View

Navigate to the **Pods** view in k9s (usually with `:pods` or `:po`):

| Action | Key | Effect |
|--------|-----|--------|
| **Tap with Mittens** | `Ctrl+M` | Tap the Service associated with this pod |

This works by extracting the pod's `app` label to determine the service name.

### Example Workflow

1. In k9s, navigate to **Services** view (press `:svc`)
2. Select the service you want to tap
3. Press `Ctrl+M` to tap it with mittens
4. Watch the pod spin up with the mitmproxy sidecar
5. Use `kubectl exec` or connect manually to interact with mitmproxy
6. When done, use `kubectl mittens off <service>` to remove mittens
7. The cleanup happens automatically

## Plugin Details

The k9s plugins are defined in `docs/k9s-plugin.yaml` and configured as follows:

- **mittens-tap**: Runs `kubectl mittens on $NAME -n $NAMESPACE`
  - Shortcut: `Ctrl+M`
  - Scope: Services view
  
- **mittens-tap-pod**: Runs `kubectl mittens on <service> -n $NAMESPACE` (extracts service from pod's `app` label)
  - Shortcut: `Ctrl+M`
  - Scope: Pods view

## Notes

- The plugins automatically use the selected service name and namespace from k9s
- All commands run in foreground mode so you can see the output
- Make sure `kubectl-mittens` is installed and in your PATH
- You can customize the shortcuts by editing `~/.k9s/plugins.yaml`

For more information on k9s plugins, see the [k9s plugins documentation](https://k9scli.io/topics/plugins/).

#!/bin/bash

# Note: pipefail is not POSIX, so we skip it here for compatibility
# The script will still work correctly without it

# Copy the config file if it exists and is readable
if [ -f /home/mitmproxy/config/config.yaml ] && [ -r /home/mitmproxy/config/config.yaml ]; then
  cp /home/mitmproxy/config/config.yaml /home/mitmproxy/.mitmproxy/config.yaml
  echo "Config file copied to /home/mitmproxy/.mitmproxy/config.yaml" >&2
else
  echo "Warning: Config file not found or not readable at /home/mitmproxy/config/config.yaml" >&2
fi

prog="${1}"
case "$prog" in
  mitmproxy)
    MITMPROXY_PATH='/home/mitmproxy/.mitmproxy'
    
    # Start a tmux session with mitmproxy to allow interactive access without requiring a TTY
    # Users can 'kubectl exec -it <pod> -- tmux attach-session -t mitmproxy' to interact
    echo "Starting mitmproxy in tmux session with confdir=${MITMPROXY_PATH}" >&2
    
    # Create a tmux session and capture any startup errors
    tmux new-session -d -s mitmproxy -c /home/mitmproxy -x 200 -y 50 \
      "mitmproxy --set confdir=${MITMPROXY_PATH} ${@:2}; bash"
    
    echo "Mitmproxy tmux session created. Waiting for it to start..." >&2
    sleep 3
    
    # Check if the session is still running
    if tmux has-session -t mitmproxy 2>/dev/null; then
      echo "Mitmproxy session is running" >&2
      echo "=== Current tmux pane content ===" >&2
      tmux capture-pane -t mitmproxy -p >&2
      echo "=== End tmux pane content ===" >&2
    else
      echo "ERROR: Mitmproxy session exited immediately" >&2
    fi
    
    # Keep the container running - sleep indefinitely
    # This allows users to attach via: kubectl exec -it <pod> -- tmux attach-session -t mitmproxy
    # or just: kubectl exec -it <pod> -- bash
    echo "Container keeping alive with sleep infinity" >&2
    sleep infinity
    ;;
  mitmdump|mitmweb)
    MITMPROXY_PATH='/home/mitmproxy/.mitmproxy'
    # For mitmdump or mitmweb, use direct execution
    echo "Starting ${prog} with confdir=${MITMPROXY_PATH}" >&2
    exec "${@}" --set "confdir=${MITMPROXY_PATH}"
    ;;
  bash|/bin/bash|sh|/bin/sh)
    # For shell commands, try to auto-attach if tmux session exists
    if tmux has-session -t mitmproxy 2>/dev/null; then
      echo "Attaching to mitmproxy session..." >&2
      exec tmux attach-session -t mitmproxy
    else
      echo "Running shell (no mitmproxy session available)" >&2
      exec "${@}"
    fi
    ;;
  "")
    # If no command specified, default to bash which will auto-attach to tmux
    if tmux has-session -t mitmproxy 2>/dev/null; then
      echo "Attaching to mitmproxy session..." >&2
      exec tmux attach-session -t mitmproxy
    else
      exec /bin/bash
    fi
    ;;
  *)
    # For any other command, just run it
    echo "Running command: ${@}" >&2
    exec "${@}"
    ;;
esac

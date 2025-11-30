#!/bin/bash

set -o pipefail

# Copy the config file if it exists and is readable
if [ -f /home/mitmproxy/config/config.yaml ] && [ -r /home/mitmproxy/config/config.yaml ]; then
  cp /home/mitmproxy/config/config.yaml /home/mitmproxy/.mitmproxy/config.yaml
  echo "Config file copied to /home/mitmproxy/.mitmproxy/config.yaml" >&2
else
  echo "Warning: Config file not found or not readable at /home/mitmproxy/config/config.yaml" >&2
fi

prog="${1}"
if [[ "${1}" == 'mitmdump' || "${1}" == 'mitmproxy' || "${1}" == 'mitmweb' ]]; then
  MITMPROXY_PATH='/home/mitmproxy/.mitmproxy'
  
  # For mitmproxy interactive terminal mode, use tmux to handle TTY requirements
  if [[ "${1}" == 'mitmproxy' ]]; then
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
    
    # Create a bashrc file that auto-attaches to tmux on interactive shell
    cat > /root/.bashrc << 'BASHRC_EOF'
# Auto-attach to mitmproxy tmux session if it exists and we're in an interactive shell
if [[ $- == *i* ]] && [ -z "$TMUX" ]; then
  if tmux has-session -t mitmproxy 2>/dev/null; then
    exec tmux attach-session -t mitmproxy
  fi
fi
BASHRC_EOF
    
    # Keep the container running - sleep indefinitely
    # This allows users to attach via: kubectl exec -it <pod> -- tmux attach-session -t mitmproxy
    # or just: kubectl exec -it <pod> -- bash
    echo "Container keeping alive with sleep infinity" >&2
    sleep infinity
  else
    # For mitmdump or mitmweb (or other commands), use direct execution
    echo "Starting ${prog} with confdir=${MITMPROXY_PATH}" >&2
    exec "${@}" --set "confdir=${MITMPROXY_PATH}"
  fi
else
  # For interactive shells, auto-attach to mitmproxy if available
  if [[ "$1" == "bash" || "$1" == "sh" || -z "$1" ]]; then
    if tmux has-session -t mitmproxy 2>/dev/null; then
      exec tmux attach-session -t mitmproxy
    fi
  fi
  echo "Running command: ${@}" >&2
  exec "${@}"
fi
